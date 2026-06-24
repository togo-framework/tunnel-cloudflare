<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/tunnel-cloudflare</h1>
  <p>Cloudflare Tunnel driver for togo tunnel — quick (*.trycloudflare.com) or named.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/tunnel-cloudflare"><img src="https://pkg.go.dev/badge/github.com/togo-framework/tunnel-cloudflare.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/tunnel-cloudflare
```
<!-- /togo-header -->

**Cloudflare Tunnel** driver for togo's [`tunnel`](https://github.com/togo-framework/tunnel)
subsystem. Wraps the `cloudflared` binary: a **quick tunnel** (no account, a
`*.trycloudflare.com` URL) by default, or a **named tunnel** when a token is set.

Requires the [`cloudflared`](https://github.com/cloudflare/cloudflared) binary on `PATH`.

## Config

| Env | Meaning |
|-----|---------|
| `TUNNEL_DRIVER` | set to `cloudflare` |
| `CLOUDFLARE_TUNNEL_TOKEN` | named-tunnel connector token (optional; omit for a quick tunnel) |
| `CLOUDFLARE_TUNNEL_HOSTNAME` | the hostname a named tunnel routes to (reported as the URL) |
| `CLOUDFLARED_BIN` | path to `cloudflared` (default: `cloudflared` on PATH) |

```go
svc, _ := tunnel.FromKernel(k)
url, _ := svc.Start(ctx, "8080")   // → https://<name>.trycloudflare.com
defer svc.Stop(ctx)
```

Quick tunnels parse the public URL from `cloudflared` output (up to 45s). Named
tunnels start the connector and report `CLOUDFLARE_TUNNEL_HOSTNAME`.

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- togo-sponsors -->
