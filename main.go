package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	flag "github.com/spf13/pflag"
	"tailscale.com/tsnet"
	"tailscale.com/types/logger"
)

func main() {
	stateDir := flag.StringP("statedir", "d", "", "Directory for storing state")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("Usage: ./funnelproxy [-d|--statedir <dir>] <hostname> <target-url>")
	}

	targetURL, err := url.Parse(flag.Arg(1))
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	s := &tsnet.Server{
		Dir:      *stateDir,
		Logf:     logger.Discard,
		Hostname: flag.Arg(0),
	}
	defer s.Close()

	ln, err := s.ListenFunnel("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("Listening on https://%v\n", s.CertDomains()[0])

	err = http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyURL := targetURL.ResolveReference(r.URL)
		proxyRequest, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
		if err != nil {
			http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
			return
		}

		proxyRequest.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(proxyRequest)
		if err != nil {
			http.Error(w, "Error executing proxy request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	}))
	log.Fatal(err)
}
