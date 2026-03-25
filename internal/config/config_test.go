package config

import "testing"

func TestIsProxyURIRecognizesHTTPAndSOCKS5(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want bool
	}{
		{name: "http", uri: "http://alice:secret@example.com:8080", want: true},
		{name: "socks5", uri: "socks5://alice:secret@example.com:1080", want: true},
		{name: "vmess", uri: "vmess://example", want: true},
		{name: "invalid", uri: "ftp://example.com", want: false},
	}

	for _, tt := range tests {
		if got := IsProxyURI(tt.uri); got != tt.want {
			t.Fatalf("%s: IsProxyURI(%q) = %v, want %v", tt.name, tt.uri, got, tt.want)
		}
	}
}

func TestParseForwardProxyURL(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "empty", value: "", wantErr: false},
		{name: "http", value: "http://127.0.0.1:7890", wantErr: false},
		{name: "https", value: "https://127.0.0.1:7890", wantErr: false},
		{name: "socks5", value: "socks5://127.0.0.1:7890", wantErr: false},
		{name: "socks5h", value: "socks5h://127.0.0.1:7890", wantErr: false},
		{name: "unsupported", value: "ftp://127.0.0.1:7890", wantErr: true},
		{name: "missing_host", value: "http://", wantErr: true},
	}

	for _, tt := range tests {
		_, err := ParseForwardProxyURL(tt.value)
		if (err != nil) != tt.wantErr {
			t.Fatalf("%s: ParseForwardProxyURL(%q) err=%v, wantErr=%v", tt.name, tt.value, err, tt.wantErr)
		}
	}
}
