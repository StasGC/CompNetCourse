package main

import (
	"fmt"
	"net"
	"time"
)

func FreePort(ip string, port int, timeout time.Duration) error {
	host := fmt.Sprintf("%s:%v", ip, port)

	_, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func main() {
	var ip string
	fmt.Println("Enter IP:")
	fmt.Scanf("%s", &ip)

	fmt.Println("Available ports:")
	for port := 0; port < 65535; port++ {
		if err := FreePort(ip, port, 5*time.Millisecond); err == nil {
			fmt.Printf("%s:%v\n", ip, port)
		}
	}
}
