package builder

import (
	"log"
	"strings"
)

// pluginNameMap maps subscription plugin names to actual executable names.
// This is necessary because subscription URIs may use different naming conventions
// than the actual plugin binaries installed in the system.
var pluginNameMap = map[string]string{
	"simple-obfs":  "obfs-local",
	"obfs":         "obfs-local",
	"v2ray-plugin": "v2ray-plugin",
	"xray-plugin":  "xray-plugin",
	"kcptun":       "kcptun",
}

// parsePluginString splits a plugin string like "simple-obfs;obfs=http;obfs-host=example.com"
// into the plugin name and its options string.
// Returns (pluginName, pluginOptions)
func parsePluginString(pluginStr string) (string, string) {
	if pluginStr == "" {
		return "", ""
	}

	// Split by semicolon - first part is plugin name, rest are options
	parts := strings.SplitN(pluginStr, ";", 2)
	pluginName := strings.TrimSpace(parts[0])

	if len(parts) > 1 {
		// Join remaining parts as options
		pluginOpts := strings.TrimSpace(parts[1])
		return pluginName, pluginOpts
	}

	return pluginName, ""
}

// normalizePluginName maps subscription plugin names to actual binary names.
// It performs case-insensitive matching and logs when mappings occur to aid
// in debugging subscription parsing issues.
//
// Parameters:
//   - pluginName: The plugin name from the subscription URI
//
// Returns:
//   - The normalized plugin binary name, or the original name if no mapping exists
func normalizePluginName(pluginName string) string {
	if pluginName == "" {
		return ""
	}

	normalized := strings.TrimSpace(pluginName)
	key := strings.ToLower(normalized)

	if mapped, ok := pluginNameMap[key]; ok {
		if mapped != normalized {
			log.Printf("🔄 Plugin name mapping: '%s' → '%s'", pluginName, mapped)
		}
		return mapped
	}

	// Unknown plugins pass through unchanged to allow flexibility
	// for future plugins or custom builds
	log.Printf("⚠️  Unknown shadowsocks plugin '%s', using as-is", pluginName)
	return normalized
}
