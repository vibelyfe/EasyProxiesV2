package builder

import (
	"testing"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func TestBuildNodeOutboundSupportsSOCKS5(t *testing.T) {
	outbound, err := buildNodeOutbound("socks-node", "socks5://demo:secret@99.144.123.135:30350", false)
	if err != nil {
		t.Fatalf("buildNodeOutbound returned error: %v", err)
	}
	if outbound.Type != C.TypeSOCKS {
		t.Fatalf("outbound type = %q, want %q", outbound.Type, C.TypeSOCKS)
	}

	opts, ok := outbound.Options.(*option.SOCKSOutboundOptions)
	if !ok {
		t.Fatalf("outbound options type = %T, want *option.SOCKSOutboundOptions", outbound.Options)
	}
	if opts.Server != "99.144.123.135" {
		t.Fatalf("server = %q, want %q", opts.Server, "99.144.123.135")
	}
	if opts.ServerPort != 30350 {
		t.Fatalf("server port = %d, want %d", opts.ServerPort, 30350)
	}
	if opts.Username != "demo" {
		t.Fatalf("username = %q, want %q", opts.Username, "demo")
	}
	if opts.Password != "secret" {
		t.Fatalf("password = %q, want %q", opts.Password, "secret")
	}
	if opts.Version != "5" {
		t.Fatalf("version = %q, want %q", opts.Version, "5")
	}
}

func TestBuildNodeOutboundSupportsHTTP(t *testing.T) {
	outbound, err := buildNodeOutbound("http-node", "http://alice:wonderland@example.com:8080/proxy", false)
	if err != nil {
		t.Fatalf("buildNodeOutbound returned error: %v", err)
	}
	if outbound.Type != C.TypeHTTP {
		t.Fatalf("outbound type = %q, want %q", outbound.Type, C.TypeHTTP)
	}

	opts, ok := outbound.Options.(*option.HTTPOutboundOptions)
	if !ok {
		t.Fatalf("outbound options type = %T, want *option.HTTPOutboundOptions", outbound.Options)
	}
	if opts.Server != "example.com" {
		t.Fatalf("server = %q, want %q", opts.Server, "example.com")
	}
	if opts.ServerPort != 8080 {
		t.Fatalf("server port = %d, want %d", opts.ServerPort, 8080)
	}
	if opts.Username != "alice" {
		t.Fatalf("username = %q, want %q", opts.Username, "alice")
	}
	if opts.Password != "wonderland" {
		t.Fatalf("password = %q, want %q", opts.Password, "wonderland")
	}
	if opts.Path != "/proxy" {
		t.Fatalf("path = %q, want %q", opts.Path, "/proxy")
	}
}

func TestApplyOutboundDetour(t *testing.T) {
	outbound, err := buildNodeOutbound("http-node", "http://alice:wonderland@example.com:8080", false)
	if err != nil {
		t.Fatalf("buildNodeOutbound returned error: %v", err)
	}

	if ok := applyOutboundDetour(&outbound, "forward-proxy"); !ok {
		t.Fatalf("applyOutboundDetour returned false")
	}

	opts, ok := outbound.Options.(*option.HTTPOutboundOptions)
	if !ok {
		t.Fatalf("outbound options type = %T, want *option.HTTPOutboundOptions", outbound.Options)
	}
	if opts.Detour != "forward-proxy" {
		t.Fatalf("detour = %q, want %q", opts.Detour, "forward-proxy")
	}
}

func TestRewriteForwardProxyForOutbound(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "socks5h_to_socks5",
			input: "socks5h://127.0.0.1:7890",
			want:  "socks5://127.0.0.1:7890",
		},
		{
			name:  "https_to_http_tls",
			input: "https://127.0.0.1:7890",
			want:  "http://127.0.0.1:7890?security=tls",
		},
	}

	for _, tt := range tests {
		got, err := rewriteForwardProxyForOutbound(tt.input)
		if err != nil {
			t.Fatalf("%s: rewriteForwardProxyForOutbound error: %v", tt.name, err)
		}
		if got != tt.want {
			t.Fatalf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}
