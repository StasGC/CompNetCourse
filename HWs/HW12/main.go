package main

import (
	"encoding/binary"
	"math/rand"
	"net"
)

func main() {
	randomIPs := GenerateIPs(4)

	network := NewNetwork("config.json", randomIPs)
	network.RIP()
}

func GenerateIPs(n int) (randomIps []string) {
	buf := make([]byte, 4)

	for i := 0; i < n; i++ {
		ip := rand.Uint32()

		binary.LittleEndian.PutUint32(buf, ip)
		randomIps = append(randomIps, net.IP(buf).String())
	}

	return randomIps
}
