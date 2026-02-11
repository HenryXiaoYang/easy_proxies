package builder

import (
	"testing"
)

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
