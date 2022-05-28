package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"math/rand"
	"net"
	"os"
)

func dnsResolve(hostname string, isIPv6 bool, dst *net.IPAddr) error {
	ipAddr := net.ParseIP(hostname)
	if isIPv6 && ipAddr.To16() != nil {
		fmt.Printf("Using the provided ipv6 address %s for tracing\n", hostname)
		dst.IP = ipAddr
	} else if !isIPv6 && ipAddr.To4() != nil {
		fmt.Printf("Using the provided ipv4 address %s for tracing\n", hostname)
		dst.IP = ipAddr
	} else {
		ips, err := net.LookupIP(hostname)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not resolve %s\n", hostname)
			return err
		}

		for _, ip := range ips {
			if isIPv6 && ip.To16() != nil {
				dst.IP = ip
				fmt.Printf("%s resolved to %s, using this ipv6 address for tracing\n", hostname, ip)
				break
			} else if !isIPv6 && ip.To4() != nil {
				dst.IP = ip
				fmt.Printf("%s resolved to %s, using this ipv4 address for tracing\n", hostname, ip)
				break
			}
		}

		if dst.IP == nil {
			return fmt.Errorf("Could not find a valid record for %s\n", hostname)
		}
	}

	return nil
}

func createICMPEcho(ICMPTypeEcho icmp.Type) icmp.Message {
	echo := icmp.Message{
		Type: ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Int(),
			Data: []byte(""),
		},
	}

	return echo
}

func GetHostByIP(addr net.Addr) (string, error) {
	hostName, err := net.LookupAddr(addr.String())
	if err != nil {
		return "", err
	}

	return hostName[0], nil
}
