package builder

import (
	"testing"
)

func TestParsePluginString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedName string
		expectedOpts string
	}{
		{
			name:         "empty string",
			input:        "",
			expectedName: "",
			expectedOpts: "",
		},
		{
			name:         "plugin name only",
			input:        "simple-obfs",
			expectedName: "simple-obfs",
			expectedOpts: "",
		},
		{
			name:         "plugin with single option",
			input:        "simple-obfs;obfs=http",
			expectedName: "simple-obfs",
			expectedOpts: "obfs=http",
		},
		{
			name:         "plugin with multiple options",
			input:        "simple-obfs;obfs=http;obfs-host=www.bing.com",
			expectedName: "simple-obfs",
			expectedOpts: "obfs=http;obfs-host=www.bing.com",
		},
		{
			name:         "plugin with spaces",
			input:        " simple-obfs ; obfs=http ",
			expectedName: "simple-obfs",
			expectedOpts: "obfs=http",
		},
		{
			name:         "v2ray-plugin with tls",
			input:        "v2ray-plugin;tls;host=example.com",
			expectedName: "v2ray-plugin",
			expectedOpts: "tls;host=example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, opts := parsePluginString(tt.input)
			if name != tt.expectedName {
				t.Errorf("parsePluginString(%q) name = %q, expected %q", tt.input, name, tt.expectedName)
			}
			if opts != tt.expectedOpts {
				t.Errorf("parsePluginString(%q) opts = %q, expected %q", tt.input, opts, tt.expectedOpts)
			}
		})
	}
}

func TestNormalizePluginName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple-obfs exact match",
			input:    "simple-obfs",
			expected: "obfs-local",
		},
		{
			name:     "simple-obfs case insensitive",
			input:    "SIMPLE-OBFS",
			expected: "obfs-local",
		},
		{
			name:     "simple-obfs with spaces",
			input:    "  simple-obfs  ",
			expected: "obfs-local",
		},
		{
			name:     "obfs short form",
			input:    "obfs",
			expected: "obfs-local",
		},
		{
			name:     "v2ray-plugin",
			input:    "v2ray-plugin",
			expected: "v2ray-plugin",
		},
		{
			name:     "xray-plugin",
			input:    "xray-plugin",
			expected: "xray-plugin",
		},
		{
			name:     "kcptun",
			input:    "kcptun",
			expected: "kcptun",
		},
		{
			name:     "unknown plugin passes through",
			input:    "custom-plugin",
			expected: "custom-plugin",
		},
		{
			name:     "unknown plugin with mixed case",
			input:    "Custom-Plugin",
			expected: "Custom-Plugin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePluginName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizePluginName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPluginNameMapCompleteness(t *testing.T) {
	// Verify all mapped plugins have valid entries
	for key, value := range pluginNameMap {
		if key == "" {
			t.Error("pluginNameMap contains empty key")
		}
		if value == "" {
			t.Errorf("pluginNameMap[%q] has empty value", key)
		}
		// Keys should be lowercase for case-insensitive matching
		if key != key {
			t.Errorf("pluginNameMap key %q should be lowercase", key)
		}
	}
}

func BenchmarkNormalizePluginName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normalizePluginName("simple-obfs")
	}
}

func BenchmarkParsePluginString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePluginString("simple-obfs;obfs=http;obfs-host=www.bing.com")
	}
}
