package services

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// GeoLocation represents geographic location data
type GeoLocation struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

// IPAPIResponse represents the response from ipapi.co
type IPAPIResponse struct {
	Country     string `json:"country_name"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Error       bool   `json:"error"`
	Reason      string `json:"reason"`
}

// GetLocationFromIP gets geographic location data from an IP address
func GetLocationFromIP(ipAddress string) (*GeoLocation, error) {
	// Handle local/private IPs
	if isLocalIP(ipAddress) {
		return &GeoLocation{
			Country: "Local",
			Region:  "Local",
			City:    "Local",
		}, nil
	}

	// Use ipapi.co free tier (1000 requests per month)
	url := fmt.Sprintf("https://ipapi.co/%s/json/", ipAddress)
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var ipData IPAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if ipData.Error {
		return nil, fmt.Errorf("API error: %s", ipData.Reason)
	}

	return &GeoLocation{
		Country: getStringOrDefault(ipData.Country, "Unknown"),
		Region:  getStringOrDefault(ipData.Region, "Unknown"),
		City:    getStringOrDefault(ipData.City, "Unknown"),
	}, nil
}

// isLocalIP checks if an IP address is local/private
func isLocalIP(ipAddress string) bool {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return true // If we can't parse it, treat as local
	}

	// Check for IPv4 private ranges
	if ip.To4() != nil {
		return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
	}

	// Check for IPv6 local addresses
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || strings.HasPrefix(ip.String(), "fe80:")
}

// getStringOrDefault returns the string if not empty, otherwise returns the default
func getStringOrDefault(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

// GetLocationFromIPFallback provides a basic fallback for geographic data
// This is used when the external API is unavailable
func GetLocationFromIPFallback(ipAddress string) *GeoLocation {
	if isLocalIP(ipAddress) {
		return &GeoLocation{
			Country: "Local",
			Region:  "Local",
			City:    "Local",
		}
	}

	return &GeoLocation{
		Country: "Unknown",
		Region:  "Unknown", 
		City:    "Unknown",
	}
} 