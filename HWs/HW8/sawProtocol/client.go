package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const host = "127.0.0.1"
const port = ":8081"
const maxBufferSize = 1024
const shard = 5

func init() {
	fmt.Println("Start client...")
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var timeout time.Duration = 1
	var fileName string = "test_data.txt"

	fmt.Println("Enter timeout(in seconds):")
	fmt.Scanln(&timeout)
	fmt.Println("Enter filename:")
	fmt.Scanln(&fileName)

	dst, err := net.ResolveUDPAddr("udp", host+port)
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, dst)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	lenData := len(data)
	if lenData > maxBufferSize {
		log.Fatal("Massage bigger then max size.")
		return
	}

	packets := make([][]byte, 0, shard)
	for i := 0; i < lenData; i += shard {
		if i+shard > lenData {
			packets = append(packets, data[i:lenData])
		} else {
			packets = append(packets, data[i:i+shard])
		}
	}

	count := 0
	current := 0
	for count < len(packets) {
		currPacket := make([]byte, len(packets[count]))
		copy(currPacket, packets[count])
		currPacket = append(currPacket, byte(current))

		if rand.Intn(10) >= 3 {
			_, err := conn.Write(currPacket)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		ackBuf := make([]byte, maxBufferSize)
		deadline := time.Now().Add(timeout * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			log.Fatal(err)
			return
		}
		n, err := conn.Read(ackBuf)
		if err != nil {
			// try to send again
			continue
		}

		if int(ackBuf[n-1]) != current {
			// try to send again
			continue
		}

		count++
		current = (current + 1) % 2
	}

	_, err = conn.Write([]byte("END"))
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Exiting UDP client!")
}
