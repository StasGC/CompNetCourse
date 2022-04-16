package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const port = ":8081"
const maxBufferSize = 1024

func init() {
	fmt.Println("Start server...")
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var fileName string = "server_test_data.txt"

	dst, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := net.ListenUDP("udp", dst)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	fileBuffer := make([]byte, maxBufferSize)
	readBuffer := make([]byte, maxBufferSize)
	current := 0
	count := 0

	for {
		n, addr, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Fatal(err)
			return
		}

		if strings.TrimSpace(string(readBuffer[0:n])) == "END" {
			break
		}

		if readBuffer[n-1] != byte(current) {
			sendAck(conn, byte((current+1)%2), addr)
			continue
		}

		for i, j := count, 0; i < count+n-1; i, j = i+1, j+1 {
			fileBuffer[i] = readBuffer[j]
		}
		count += n - 1

		sendAck(conn, byte(current), addr)
		current = (current + 1) % 2
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = file.Write(fileBuffer[0:count])
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Exiting UDP server!")
}

func sendAck(conn *net.UDPConn, currPacket byte, addr *net.UDPAddr) {
	ackBuffer := []byte("ACK")
	ackBuffer = append(ackBuffer, currPacket)

	if rand.Intn(10) >= 3 {
		_, err := conn.WriteToUDP(ackBuffer, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
