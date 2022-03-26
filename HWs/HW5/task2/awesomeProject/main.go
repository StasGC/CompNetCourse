package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	readString, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	readString = strings.TrimSpace(readString)

	fmt.Print(readString)
	fmt.Print(readString == "hello")

}
