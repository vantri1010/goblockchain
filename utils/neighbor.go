package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("%s %v\n", target, err)
		return false
	}
	return true
}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
	// Construct the current node's address in "host:port" format (e.g., "192.168.1.100:8080")
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	// Use regex PATTERN to parse the IP address (e.g., "192.168.1.100" into "192.168.1." and "100")
	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		// Return nil if the host format is invalid (regex doesn't match)
		return nil
	}
	// Extract the IP prefix (e.g., "192.168.1.") and the last octet (e.g., "100")
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])

	// Initialize an empty slice to store discovered neighbor addresses
	neighbors := make([]string, 0)

	// Iterate over the specified port range (startPort to endPort)
	for port := startPort; port <= endPort; port += 1 {
		// Iterate over the specified IP offset range (startIp to endIp)
		for ip := startIp; ip <= endIp; ip += 1 {
			// Construct a potential neighbor's IP by appending the adjusted last octet
			// E.g., if prefixHost="192.168.1.", lastIp=100, ip=1, then guessHost="192.168.1.101"
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			// Form the full address (e.g., "192.168.1.101:8000")
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			// Add to neighbors if the address is not the current node and the host is reachable
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	// Return the list of discovered neighbor addresses
	return neighbors
}
func GetHost() string {
	hostname, err := os.Hostname()
	fmt.Println(hostname)
	if err != nil {
		return "127.0.0.1"
	}
	address, err := net.LookupHost(hostname)
	fmt.Println(address)
	if err != nil {
		return "127.0.0.1"
	}
	return address[0]
}
