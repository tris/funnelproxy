# funnelproxy

## Overview

This is a quick & dirty HTTP proxy to expose a URL on its own hostname through
Tailscale Funnel.

This can be thought of as equivalent to running the following on a Tailscale
node:

```bash
sudo tailscale serve https / https://target-url
sudo tailscale funnel 443 on
```

The difference is, funnelproxy creates its own hostname for the served URL.
This can be useful in cases where you need to serve more than one webapp, and
those webapps don't work with relative paths.

## Prerequisites

1. Follow the Tailscale Funnel [setup instructions](https://tailscale.com/kb/1223/tailscale-funnel/#setup)
(if you have not used Funnel before)
2. Generate an auth key at [Admin console > Settings > Keys](https://login.tailscale.com/admin/settings/keys)

## Usage

```bash
go install github.com/tris/funnelproxy@latest
TS_AUTHKEY=tskey-auth-xxx funnelproxy fun http://target-url
```

This will create `https://fun.<tailnet>.ts.net`, forwarding all requests to
`http://target-url`.

### Docker

```bash
docker volume create fp_fun
TS_AUTHKEY=tskey-auth-xxx docker run --rm \
  --name fp_fun \
  -v fp_fun:/state \
  -e TS_AUTHKEY \
  ghcr.io/tris/funnelproxy -- -d /state fun http://target-url
```

### Docker Compose

```yaml
version: '3'

volumes:
  fp_fun: {}

services:
  fp_fun:
    image: ghcr.io/tris/funnelproxy
    volumes:
    - fp_fun:/state
    environment:
    - TS_AUTHKEY=tskey-auth-xxx
    command:
    - '--statedir=/state'
    - 'fun'
    - 'http://target-url'
    restart: unless-stopped
```

## See also

* https://tailscale.com/blog/tailscale-funnel-beta/
* https://tailscale.dev/blog/embedded-funnel
* https://tailscale.com/kb/1085/auth-keys/
