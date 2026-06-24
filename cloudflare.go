// Package cloudflare is a togo tunnel driver for Cloudflare Tunnel. It wraps the
// `cloudflared` binary: a quick tunnel (no account, *.trycloudflare.com) by
// default, or a named tunnel when CLOUDFLARE_TUNNEL_TOKEN is set.
//
// Install: `togo install togo-framework/tunnel-cloudflare`, set TUNNEL_DRIVER=cloudflare.
// Requires the `cloudflared` binary on PATH (https://github.com/cloudflare/cloudflared).
package cloudflare

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"

	"github.com/togo-framework/togo"
	"github.com/togo-framework/tunnel"
)

func init() {
	tunnel.RegisterDriver("cloudflare", func(k *togo.Kernel) (tunnel.Tunnel, error) {
		return &driver{
			bin:   envOr("CLOUDFLARED_BIN", "cloudflared"),
			token: os.Getenv("CLOUDFLARE_TUNNEL_TOKEN"),
			host:  os.Getenv("CLOUDFLARE_TUNNEL_HOSTNAME"),
		}, nil
	})
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

type driver struct {
	bin, token, host string

	mu  sync.Mutex
	cmd *exec.Cmd
	url string
}

// quickURLRe matches the public URL cloudflared prints for a quick tunnel.
var quickURLRe = regexp.MustCompile(`https://[a-z0-9-]+\.trycloudflare\.com`)

// parseQuickURL extracts the *.trycloudflare.com URL from a line of cloudflared
// output, or "" if none.
func parseQuickURL(line string) string { return quickURLRe.FindString(line) }

func (d *driver) Start(ctx context.Context, addr string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.cmd != nil {
		return d.url, nil
	}
	local := "http://" + tunnel.NormalizeAddr(addr)

	var cmd *exec.Cmd
	named := d.token != ""
	if named {
		// Named tunnel: routing to a hostname is configured in Cloudflare; the
		// token wires the connector. We can't discover the hostname from output,
		// so the caller supplies CLOUDFLARE_TUNNEL_HOSTNAME.
		cmd = exec.Command(d.bin, "tunnel", "run", "--token", d.token)
	} else {
		cmd = exec.Command(d.bin, "tunnel", "--no-autoupdate", "--url", local)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	cmd.Stdout = io.Discard
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("tunnel-cloudflare: start %s: %w (is cloudflared installed?)", d.bin, err)
	}
	d.cmd = cmd

	if named {
		d.url = d.host
		go drain(stderr)
		if d.url == "" {
			return "", fmt.Errorf("tunnel-cloudflare: named tunnel started; set CLOUDFLARE_TUNNEL_HOSTNAME to report its URL")
		}
		return d.url, nil
	}

	// Quick tunnel: scan stderr for the trycloudflare URL.
	urlCh := make(chan string, 1)
	go func() {
		sc := bufio.NewScanner(stderr)
		for sc.Scan() {
			if u := parseQuickURL(sc.Text()); u != "" {
				select {
				case urlCh <- u:
				default:
				}
			}
		}
	}()

	timeout := 45 * time.Second
	select {
	case <-ctx.Done():
		_ = d.stopLocked()
		return "", ctx.Err()
	case u := <-urlCh:
		d.url = u
		return u, nil
	case <-time.After(timeout):
		_ = d.stopLocked()
		return "", fmt.Errorf("tunnel-cloudflare: timed out after %s waiting for the trycloudflare URL", timeout)
	}
}

func drain(r io.Reader) { _, _ = io.Copy(io.Discard, r) }

func (d *driver) Stop(context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.stopLocked()
}

func (d *driver) stopLocked() error {
	if d.cmd == nil || d.cmd.Process == nil {
		return nil
	}
	err := d.cmd.Process.Kill()
	_ = d.cmd.Wait()
	d.cmd = nil
	d.url = ""
	return err
}

func (d *driver) Status(context.Context) (tunnel.Status, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return tunnel.Status{Running: d.cmd != nil, URL: d.url}, nil
}
