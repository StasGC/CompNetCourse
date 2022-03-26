package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Client for remote commands launching...")

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	consoleReader := bufio.NewReader(os.Stdin)
	//connReader := bufio.NewReader(conn)

	for {
		fmt.Println("\n" + "Call new command.")
		fmt.Println("For close connection type 'close'.")
		fmt.Println("Text command to send:")

		command, err := consoleReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		command = strings.TrimSpace(command)

		if command == "close" {
			fmt.Println("Connection closed.")
			return
		}

		fmt.Fprintf(conn, "/c "+command+"\n")

		//for {
		//	message, err := connReader.ReadString('\n')
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//	message = strings.TrimSpace(message)
		//	if strings.Index(message, endTokenClient) != -1 {
		//		break
		//	}
		//
		//	fmt.Println(message)
		//}
	}
}
