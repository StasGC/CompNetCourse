package main

import (
	"fmt"
	"net"
	"strings"
)

func handleMask(mask net.IPMask) string {
	stringMusk := mask.String()
	newMusk := make([]string, 4)

	for i := 0; i < 8; i += 2 {
		if stringMusk[i:i+2] == "ff" {
			newMusk = append(newMusk, "255")
			newMusk = append(newMusk, ".")
		}
		if stringMusk[i:i+2] == "00" {
			newMusk = append(newMusk, "0")
			newMusk = append(newMusk, ".")
		}
	}

	newMusk = newMusk[0 : len(newMusk)-1]

	return strings.Join(newMusk, "")

}

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		if aspnet, ok := address.(*net.IPNet); ok && !aspnet.IP.IsLoopback() {
			if aspnet.IP.To4() != nil {
				fmt.Println("IP:", aspnet.IP.String())
				fmt.Println("Mask:", handleMask(aspnet.IP.DefaultMask()))
				fmt.Println("")
			}
		}
	}
}
