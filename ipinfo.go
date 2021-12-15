// Package ipinfo provides info on IP address location
// using the ipinfo.io service.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

// IPInfo describes a particular IP address.
type IPInfo struct {
	// IP holds the described IP address.
	IP string
	// Hostname holds a DNS name associated with the IP address.
	Hostname string
	// City holds the city of the ISP location.
	City string
	// Region holds the region of the ISP location.
	Region string
	// Country holds the two-letter country code.
	Country string
	// Loc holds the latitude and longitude of the
	// ISP location as a comma-separated northing, easting
	// pair of floating point numbers.
	Loc string
	// Org describes the organization that is
	// responsible for the IP address.
	Org string
	// Timezone holds the timezone of the ISP location.
	Timezone string
	// Postal holds the post code or zip code region of the ISP location.
	Postal string
}

// MyIP provides information about the public IP address of the client.
func MyIP() (*IPInfo, error) {
	return ForeignIP("")
}

// ForeignIP provides information about the given IP address,
// which should be in dotted-quad form.
// Example:
// http://ipinfo.io/json
// http://ipinfo.io/67.148.93.12/json
func ForeignIP(ip string) (*IPInfo, error) {
	var ipinfo IPInfo
	if ip == "" || IsPrivateIP(net.ParseIP(ip)) {
		if err := GetV4Redis(RIPInfo("127.0.0.1"), &ipinfo); err == nil {
			return &ipinfo, nil
		}
		ip = ""
	}
	if ip != "" {
		ip += "/" + ip
	}
	response, err := http.Get("http://ipinfo.io" + ip + "/json")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(contents, &ipinfo); err != nil {
		return nil, err
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	SetV2Redis(RIPInfo(ip), &ipinfo, ExpireTimeIPInfo)
	return &ipinfo, nil
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

// IsPrivateIP returns true if ip is private
func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
