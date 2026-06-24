package cloudflare

import (
	"testing"

	"github.com/togo-framework/tunnel"
)

func TestParseQuickURL(t *testing.T) {
	cases := map[string]string{
		`2024-01-02T03:04:05Z INF +-----+ |  https://happy-tree-1234.trycloudflare.com  | +-----+`: "https://happy-tree-1234.trycloudflare.com",
		`INF Your quick Tunnel has been created! Visit it at https://red-fox-9.trycloudflare.com`:    "https://red-fox-9.trycloudflare.com",
		`INF Registered tunnel connection connIndex=0`:                                               "",
		``: "",
	}
	for line, want := range cases {
		if got := parseQuickURL(line); got != want {
			t.Errorf("parseQuickURL(%q) = %q, want %q", line, got, want)
		}
	}
}

func TestDriverRegistered(t *testing.T) {
	found := false
	for _, n := range tunnel.Drivers() {
		if n == "cloudflare" {
			found = true
		}
	}
	if !found {
		t.Fatal("cloudflare driver not registered on tunnel base")
	}
}

func TestEnvOr(t *testing.T) {
	if envOr("CLOUDFLARED_BIN_UNSET_XYZ", "cloudflared") != "cloudflared" {
		t.Error("envOr default failed")
	}
}
