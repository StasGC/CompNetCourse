package main

import (
	"flag"
	"fmt"
)

func main() {
	remoteAddress := flag.String("remote", "", "Remote host (can be an ip or hostname)")
	localAddress := flag.String("local", "", "local bind address")
	numMessages := flag.Int("num-messages", 3, "The number of messages that are sent to each router")
	ttl := flag.Int("ttl", 20, "ttl to use (default: 20)")
	timeout := flag.Int("timeout", 2, "timeout to wait for an ICMP answer")
	flag.Parse()

	if *remoteAddress == "" {
		fmt.Printf("Remote address is empty\n")
		return
	}

	t := NewTracerouteSession(*localAddress, *remoteAddress, *ttl, *timeout, *numMessages)
	t.runTrace()
}
