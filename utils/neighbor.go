package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 3*time.Second) // Increased timeout to 3 seconds
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", target, err)
		return false
	}
	fmt.Printf("Successfully connected to %s\n", target)
	return true
}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
	// Construct the current node's address in "host:port" format (e.g., "192.168.1.100:8080")
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	// Use regex PATTERN to parse the IP address (e.g., "192.168.1.100" into "192.168.1." and "100")
	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		fmt.Printf("Invalid host format: %s\n", myHost)
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
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			// Form the full address (e.g., "192.168.1.101:8000")
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			// Add to neighbors if the address is not the current node and the host is reachable
			fmt.Printf("Trying to connect to %s\n", guessTarget)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	fmt.Println("Find neighbors: ", neighbors)
	// Return the list of discovered neighbor addresses
	return neighbors
}

func GetHost() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error getting interfaces: %v\n", err)
		return "127.0.0.1"
	}
	for _, iface := range interfaces {
		// Skip down interfaces and loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				// Prefer IPv4 and non-loopback IPs
				if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() {
					fmt.Printf("Found IP: %s on interface %s\n", ipnet.IP.String(), iface.Name)
					return ipnet.IP.String()
				}
			}
		}
	}
	fmt.Println("No suitable IP found, falling back to 127.0.0.1")
	return "127.0.0.1"
}
