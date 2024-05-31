package corsproxy

import (
	"testing"
)

func TestWildcardMatch(t *testing.T) {
	tests := []struct {
		wildcard wildcard
		target   string
		expected bool
	}{
		{wildcard{"http://example.com/", ""}, "http://example.com/", true},
		{wildcard{"http://example.com/", ""}, "http://example.com/sub", true},
		{wildcard{"http://example.com", ""}, "http://example.net", false},
		{wildcard{"http://", ".example.com"}, "http://example.com", false},
		{wildcard{"http://", ".example.com"}, "http://sub.example.com", true},
		{wildcard{"http://example.com/sub1/", "/sub3"}, "http://example.com/sub1/sub2/sub3", true},
		{wildcard{"http://example.com/sub1/", "/sub3"}, "http://example.com/sub1/sub2/sub3/sub4", true},
		{wildcard{"http://example.com/sub1/", "/sub3"}, "http://example.com/sub1/sub3", false},
	}

	for _, test := range tests {
		result := test.wildcard.match(test.target)
		if result != test.expected {
			t.Errorf("wildcard.match(%q) = %v; want %v", test.target, result, test.expected)
		}
	}
}

func TestIsPrivateAddr(t *testing.T) {
	tests := []struct {
		host    string
		private bool
		parsed  bool
	}{
		{"127.0.0.1", true, true},
		{"192.168.0.1", true, true},
		{"10.0.0.1", true, true},
		{"172.16.0.1", true, true},
		{"8.8.8.8", false, true},
		{"example.com", false, false},
		{"127.0.0.1:8080", true, true},
		{"192.168.0.1:8080", true, true},
		{"8.8.8.8:8080", false, true},
	}

	for _, test := range tests {
		private, parsed := isPrivateAddr(test.host)
		if private != test.private || parsed != test.parsed {
			t.Errorf("isPrivateAddr(%q) = (%v, %v); want (%v, %v)", test.host, private, parsed, test.private, test.parsed)
		}
	}
}

func TestStripURLQuery(t *testing.T) {
	tests := []struct {
		rawURL    string
		expected  string
		expectErr bool
	}{
		{"http://example.com/path?query=1", "http://example.com/path", false},
		{"http://example.com/path", "http://example.com/path", false},
		{"http://example.com", "http://example.com", false},
		{"http://EXAMPLE.COM/path?query=1", "http://example.com/path", false},
		{"not a url", "not%20a%20url", false},
	}

	for _, test := range tests {
		result, err := StripURLQuery(test.rawURL)
		if (err != nil) != test.expectErr {
			t.Errorf("StripURLQuery(%q) error = %v; wantErr %v", test.rawURL, err, test.expectErr)
		}
		if result != test.expected {
			t.Errorf("StripURLQuery(%q) = %q; want %q", test.rawURL, result, test.expected)
		}
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		rawURL    string
		expected  string
		expectErr bool
	}{
		{"http://example.com/path", "http://example.com/path", false},
		{"http://EXAMPLE.COM/path", "http://example.com/path", false},
		{"not a url", "not%20a%20url", false},
	}

	for _, test := range tests {
		result, err := NormalizeURL(test.rawURL)
		if (err != nil) != test.expectErr {
			t.Errorf("NormalizeURL(%q) error = %v; wantErr %v", test.rawURL, err, test.expectErr)
		}
		if result != test.expected {
			t.Errorf("NormalizeURL(%q) = %q; want %q", test.rawURL, result, test.expected)
		}
	}
}

func TestNormalizeParseURL(t *testing.T) {
	tests := []struct {
		rawURL    string
		expected  string
		expectErr bool
	}{
		{"http://example.com/path", "example.com", false},
		{"http://EXAMPLE.COM/path", "example.com", false},
		{"not a url", "", false},
	}

	for _, test := range tests {
		result, err := NormalizeParseURL(test.rawURL)
		if (err != nil) != test.expectErr {
			t.Errorf("NormalizeParseURL(%q) error = %v; wantErr %v", test.rawURL, err, test.expectErr)
		}
		if err == nil && result.Host != test.expected {
			t.Errorf("NormalizeParseURL(%q) = %q; want %q", test.rawURL, result.Host, test.expected)
		}
	}
}
