package middleware

import (
	"net"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// IPFilter represents IP filtering configuration
type IPFilter struct {
	Whitelist []string // List of allowed IPs/CIDR ranges
	Blacklist []string // List of blocked IPs/CIDR ranges
	Mode      string   // "whitelist", "blacklist", or "both"
}

// LoadIPFilterFromEnv loads IP filter configuration from environment variables
func LoadIPFilterFromEnv() *IPFilter {
	filter := &IPFilter{
		Mode: os.Getenv("IP_FILTER_MODE"), // whitelist, blacklist, or both
	}

	// Load whitelist from environment
	if whitelist := os.Getenv("IP_WHITELIST"); whitelist != "" {
		filter.Whitelist = strings.Split(whitelist, ",")
		for i, ip := range filter.Whitelist {
			filter.Whitelist[i] = strings.TrimSpace(ip)
		}
	}

	// Load blacklist from environment
	if blacklist := os.Getenv("IP_BLACKLIST"); blacklist != "" {
		filter.Blacklist = strings.Split(blacklist, ",")
		for i, ip := range filter.Blacklist {
			filter.Blacklist[i] = strings.TrimSpace(ip)
		}
	}

	return filter
}

// isIPInRange checks if an IP address is in a given CIDR range or matches exactly
func isIPInRange(clientIP, rangeOrIP string) bool {
	// First try exact IP match
	if clientIP == rangeOrIP {
		return true
	}

	// Try CIDR range match
	_, network, err := net.ParseCIDR(rangeOrIP)
	if err != nil {
		// If not a valid CIDR, treat as single IP
		return clientIP == rangeOrIP
	}

	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	return network.Contains(ip)
}

// isIPInList checks if an IP address is in a list of IPs/CIDR ranges
func isIPInList(clientIP string, ipList []string) bool {
	for _, rangeOrIP := range ipList {
		if isIPInRange(clientIP, rangeOrIP) {
			return true
		}
	}
	return false
}

// IPFilterMiddleware creates middleware for IP filtering
func IPFilterMiddleware(filter *IPFilter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip filtering if no mode is set
		if filter.Mode == "" {
			return c.Next()
		}

		// Get client IP using the same method as analytics middleware
		clientIP := GetClientIP(c)

		switch filter.Mode {
		case "whitelist":
			if len(filter.Whitelist) > 0 && !isIPInList(clientIP, filter.Whitelist) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied: IP address not whitelisted",
				})
			}

		case "blacklist":
			if len(filter.Blacklist) > 0 && isIPInList(clientIP, filter.Blacklist) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied: IP address blacklisted",
				})
			}

		case "both":
			// Check blacklist first
			if len(filter.Blacklist) > 0 && isIPInList(clientIP, filter.Blacklist) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied: IP address blacklisted",
				})
			}
			// Then check whitelist (if whitelist exists, IP must be in it)
			if len(filter.Whitelist) > 0 && !isIPInList(clientIP, filter.Whitelist) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied: IP address not whitelisted",
				})
			}
		}

		return c.Next()
	}
}

// CreateAdminIPFilter creates an IP filter specifically for admin endpoints
func CreateAdminIPFilter() *IPFilter {
	return &IPFilter{
		Mode:      os.Getenv("ADMIN_IP_FILTER_MODE"), // separate config for admin endpoints
		Whitelist: getAdminWhitelist(),
		Blacklist: getAdminBlacklist(),
	}
}

func getAdminWhitelist() []string {
	if whitelist := os.Getenv("ADMIN_IP_WHITELIST"); whitelist != "" {
		ips := strings.Split(whitelist, ",")
		for i, ip := range ips {
			ips[i] = strings.TrimSpace(ip)
		}
		return ips
	}
	return nil
}

func getAdminBlacklist() []string {
	if blacklist := os.Getenv("ADMIN_IP_BLACKLIST"); blacklist != "" {
		ips := strings.Split(blacklist, ",")
		for i, ip := range ips {
			ips[i] = strings.TrimSpace(ip)
		}
		return ips
	}
	return nil
} 