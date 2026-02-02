package geoip

import (
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

// Region codes
const (
	RegionJP    = "jp"
	RegionKR    = "kr"
	RegionUS    = "us"
	RegionHK    = "hk"
	RegionTW    = "tw"
	RegionOther = "other"
)

// RegionInfo contains region details
type RegionInfo struct {
	Code    string // "jp", "kr", "us", "hk", "tw", "other"
	Country string // Full country name
	ISOCode string // ISO country code
}

// Lookup provides GeoIP lookup functionality
type Lookup struct {
	db   *geoip2.Reader
	mu   sync.RWMutex
	path string
}

// New creates a new GeoIP lookup instance
func New(dbPath string) (*Lookup, error) {
	if dbPath == "" {
		return &Lookup{}, nil
	}
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Lookup{db: db, path: dbPath}, nil
}

// Close closes the GeoIP database
func (l *Lookup) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.db != nil {
		return l.db.Close()
	}
	return nil
}

// IsEnabled returns true if GeoIP lookup is available
func (l *Lookup) IsEnabled() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.db != nil
}

// LookupIP returns region info for an IP address
func (l *Lookup) LookupIP(ipStr string) RegionInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.db == nil {
		return RegionInfo{Code: RegionOther, Country: "Unknown", ISOCode: ""}
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return RegionInfo{Code: RegionOther, Country: "Unknown", ISOCode: ""}
	}

	record, err := l.db.Country(ip)
	if err != nil {
		return RegionInfo{Code: RegionOther, Country: "Unknown", ISOCode: ""}
	}

	isoCode := record.Country.IsoCode
	country := record.Country.Names["en"]
	if country == "" {
		country = isoCode
	}

	return RegionInfo{
		Code:    isoCodeToRegion(isoCode),
		Country: country,
		ISOCode: isoCode,
	}
}

// LookupURI extracts server from URI and returns region info
func (l *Lookup) LookupURI(uri string) RegionInfo {
	host := extractHostFromURI(uri)
	if host == "" {
		return RegionInfo{Code: RegionOther, Country: "Unknown", ISOCode: ""}
	}

	// Resolve hostname to IP if needed
	ip := net.ParseIP(host)
	if ip == nil {
		// It's a hostname, try to resolve
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return RegionInfo{Code: RegionOther, Country: "Unknown", ISOCode: ""}
		}
		host = ips[0].String()
	}

	return l.LookupIP(host)
}

// extractHostFromURI extracts the host/IP from various proxy URI formats
func extractHostFromURI(uri string) string {
	// Handle different URI schemes
	lowerURI := strings.ToLower(uri)

	// vmess:// vless:// trojan:// ss:// ssr:// hysteria:// hysteria2:// hy2://
	if strings.HasPrefix(lowerURI, "vmess://") ||
		strings.HasPrefix(lowerURI, "vless://") ||
		strings.HasPrefix(lowerURI, "trojan://") ||
		strings.HasPrefix(lowerURI, "hysteria://") ||
		strings.HasPrefix(lowerURI, "hysteria2://") ||
		strings.HasPrefix(lowerURI, "hy2://") {
		// Standard URL format: scheme://user@host:port?params#fragment
		parsed, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		return parsed.Hostname()
	}

	if strings.HasPrefix(lowerURI, "ss://") {
		// Shadowsocks format: ss://base64(method:password)@host:port#name
		// or ss://base64@host:port#name
		return extractSSHost(uri)
	}

	if strings.HasPrefix(lowerURI, "ssr://") {
		// SSR format is base64 encoded
		return extractSSRHost(uri)
	}

	return ""
}

func extractSSHost(uri string) string {
	// Remove ss:// prefix
	content := strings.TrimPrefix(uri, "ss://")

	// Remove fragment (#name)
	if idx := strings.Index(content, "#"); idx != -1 {
		content = content[:idx]
	}

	// Check if it's the new format: base64@host:port
	if atIdx := strings.LastIndex(content, "@"); atIdx != -1 {
		hostPort := content[atIdx+1:]
		if colonIdx := strings.LastIndex(hostPort, ":"); colonIdx != -1 {
			return hostPort[:colonIdx]
		}
		return hostPort
	}

	// Old format: entire content is base64
	return ""
}

func extractSSRHost(uri string) string {
	// SSR is complex, skip for now - will be marked as "other"
	return ""
}

// isoCodeToRegion maps ISO country codes to our region codes
func isoCodeToRegion(isoCode string) string {
	switch strings.ToUpper(isoCode) {
	case "JP":
		return RegionJP
	case "KR":
		return RegionKR
	case "US":
		return RegionUS
	case "HK":
		return RegionHK
	case "TW":
		return RegionTW
	default:
		return RegionOther
	}
}

// AllRegions returns all supported region codes
func AllRegions() []string {
	return []string{RegionJP, RegionKR, RegionUS, RegionHK, RegionTW, RegionOther}
}

// RegionName returns the display name for a region code
func RegionName(code string) string {
	switch code {
	case RegionJP:
		return "Japan"
	case RegionKR:
		return "Korea"
	case RegionUS:
		return "USA"
	case RegionHK:
		return "Hong Kong"
	case RegionTW:
		return "Taiwan"
	case RegionOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// RegionEmoji returns the flag emoji for a region code
func RegionEmoji(code string) string {
	switch code {
	case RegionJP:
		return "🇯🇵"
	case RegionKR:
		return "🇰🇷"
	case RegionUS:
		return "🇺🇸"
	case RegionHK:
		return "🇭🇰"
	case RegionTW:
		return "🇹🇼"
	case RegionOther:
		return "🌍"
	default:
		return "❓"
	}
}
