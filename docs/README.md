# tunnel-cloudflare — docs

**Cloudflare Tunnel.** Expose locally via `cloudflared` — quick or named tunnel.

## Install

```bash
togo install togo-framework/tunnel-cloudflare
```

Registers on the [`tunnel`](https://github.com/togo-framework/tunnel) base; select it with **tunnel.provider in togo.yaml (or TUNNEL_DRIVER)**, then use **`togo tunnel`**.

## Interface

`Tunnel` — `Start(ctx, addr) -> publicURL`, `Stop`, `Status`.

## Configuration

| Env var | Description |
|---|---|
| `CLOUDFLARE_TUNNEL_TOKEN` | Cloudflare named-tunnel token. Optional — without it a quick `*.trycloudflare.com` tunnel is used. |
| `CLOUDFLARE_TUNNEL_HOSTNAME` | Hostname for a named tunnel. Optional. |

## Usage & notes

Requires `cloudflared` installed. Without a token it runs a **quick tunnel** and parses the `*.trycloudflare.com` URL; with `CLOUDFLARE_TUNNEL_TOKEN` it runs a **named tunnel**.

## Example

```bash
togo tunnel:start --provider cloudflare
```

## Links

- [cloudflared](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/tunnel-cloudflare)
