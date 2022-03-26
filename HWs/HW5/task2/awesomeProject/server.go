package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"net"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Launching server...")

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Listen Error:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			return
		}

		go func(conn net.Conn) {
			defer conn.Close()

			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Reader Error:", err)
				return
			}

			message = strings.TrimRight(message, "\n")
			result := WinCmdExe(message)
			fmt.Println(result)

			fmt.Println("\n" + "Waiting next command...")
		}(conn)
	}
}

func WinCmdExe(strCommand string) string {
	argsCommand := strings.Split(strCommand, " ")
	cmd := exec.Command("cmd", argsCommand...)

	stdout, _ := cmd.Output()
	d := charmap.CodePage866.NewDecoder()
	decodeOut, _ := d.Bytes(stdout)

	return string(decodeOut)
}
