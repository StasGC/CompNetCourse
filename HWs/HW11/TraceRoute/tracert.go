package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)

func (tr *TraceRouteSession) traceRoute() {
	var dst net.IPAddr

	err := dnsResolve(tr.RemoteAddr, false, &dst)
	if err != nil {
		os.Exit(1)
	}

	icmpSock, err := net.ListenPacket("ip4:icmp", tr.LocalAddr)
	if err != nil {
		fmt.Printf("Could not set a listening ICMP socket: %s\n", err)
		return
	}
	defer icmpSock.Close()

	ipv4Sock := ipv4.NewPacketConn(icmpSock)
	defer ipv4Sock.Close()

	if err := ipv4Sock.SetControlMessage(ipv4.FlagTTL|ipv4.FlagDst|ipv4.FlagInterface|ipv4.FlagSrc, true); err != nil {
		fmt.Printf("Could not set options on the ipv4 socket: %s\n", err)
		return
	}

	icmpEcho := createICMPEcho(ipv4.ICMPTypeEcho)

	readBuf := make([]byte, 1500)

Loop:
	for i := 1; i < tr.MaxTTL; i++ {
		icmpEcho.Body.(*icmp.Echo).Seq = i

		writeBuffer, err := icmpEcho.Marshal(nil)

		if err != nil {
			fmt.Printf("Could not serialize the ICMP echo request: %s\n", err)
			os.Exit(1)
		}

		if err := ipv4Sock.SetTTL(i); err != nil {
			fmt.Printf("Could not set the TTL field: %s\n", err)
			os.Exit(1)
		}

		var currentRTTs []time.Duration
		var icmpAnswers []*icmp.Message
		var hostName string
		for j := 0; j < tr.NumMessages; j++ {
			now := time.Now()

			if _, err := ipv4Sock.WriteTo(writeBuffer, nil, &dst); err != nil {
				fmt.Printf("Could not send the ICMP packet: %s\n", err)
				os.Exit(1)
			}

			if err := ipv4Sock.SetReadDeadline(time.Now().Add(tr.Timeout)); err != nil {
				fmt.Printf("Could not set the read timeout on the ipv4 socket: %s\n", err)
				os.Exit(1)
			}

			readBytes, _, hopNode, err := ipv4Sock.ReadFrom(readBuf)
			hostName, err = GetHostByIP(hopNode)
			if err != nil {
				fmt.Printf("Could not parse IP Addrs: %s\n", err)
				os.Exit(1)
			}

			if err != nil {
				fmt.Printf("%d %20s\n", i, "*")
			} else {
				icmpAnswer, err := icmp.ParseMessage(1, readBuf[:readBytes])

				if err != nil {
					fmt.Printf("Could not parse the ICMP packet from: %s\n", hopNode.String())
					os.Exit(1)
				}

				currentRTT := time.Since(now)

				currentRTTs = append(currentRTTs, currentRTT)
				icmpAnswers = append(icmpAnswers, icmpAnswer)
			}
		}

		fmt.Printf("%d	%-20s\t", i, hostName)
		for _, RTT := range currentRTTs {
			fmt.Printf("%-10v", RTT)
		}
		fmt.Println()

		for i, answer := range icmpAnswers {
			if answer.Type == ipv4.ICMPTypeTimeExceeded {
				tr.hopNum++
				break
			} else if answer.Type == ipv4.ICMPTypeEchoReply {
				fmt.Printf("%d   %20s   %20s\n", i, hostName, currentRTTs[i])
				break Loop
			}
		}
	}
}
