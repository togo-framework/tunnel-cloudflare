---
name: tunnel-cloudflare
description: Expose a local togo app via Cloudflare Tunnel — set TUNNEL_DRIVER=cloudflare and call tunnel.Start (quick *.trycloudflare.com or a named tunnel)
---

# togo tunnel-cloudflare

Cloudflare Tunnel driver for the togo `tunnel` subsystem. Wraps the `cloudflared`
binary.

## Setup

```bash
togo install togo-framework/tunnel
togo install togo-framework/tunnel-cloudflare
```

Install [`cloudflared`](https://github.com/cloudflare/cloudflared), then `.env`:

```bash
TUNNEL_DRIVER=cloudflare
# Quick tunnel (no account) is the default → https://<random>.trycloudflare.com
# Named tunnel (your domain):
# CLOUDFLARE_TUNNEL_TOKEN=...           # connector token from the dashboard
# CLOUDFLARE_TUNNEL_HOSTNAME=app.example.com
# CLOUDFLARED_BIN=cloudflared           # override binary path
```

## Use

```go
import (
	_ "github.com/togo-framework/tunnel"
	_ "github.com/togo-framework/tunnel-cloudflare"
	"github.com/togo-framework/tunnel"
)

if tn, ok := tunnel.FromKernel(k); ok {
	url, _ := tn.Start(ctx, "8080") // https://<...>.trycloudflare.com
	defer tn.Stop(ctx)
}
```

## Notes
- Quick tunnels need no Cloudflare account — great for webhook testing.
- For a named tunnel set `CLOUDFLARE_TUNNEL_TOKEN`; the URL is your configured
  `CLOUDFLARE_TUNNEL_HOSTNAME`.
- The tunnel stops when the process stops (`tn.Stop`).
