package main

import (
	"fmt"
	"time"
)

type TraceRouteSession struct {
	MaxTTL               int
	Timeout              time.Duration
	NumMessages          int
	LocalAddr            string
	RemoteAddr           string
	latencyEstimation    time.Duration
	maxLatencyBetweenHop time.Duration
	latencyBetweenHop    time.Duration
	hopOne               int
	hopTwo               int
	hopNum               int
}

func NewTracerouteSession(
	localAddr string,
	remoteAddr string,
	maxTTL int,
	timeout int,
	numMessages int,
) *TraceRouteSession {
	if localAddr == "" {
		localAddr = "0.0.0.0"
	}

	return &TraceRouteSession{
		MaxTTL:      maxTTL,
		Timeout:     time.Duration(timeout) * time.Second,
		NumMessages: numMessages,
		LocalAddr:   localAddr,
		RemoteAddr:  remoteAddr,
		hopNum:      0,
		hopOne:      -1,
		hopTwo:      -1,
	}
}

func (tr *TraceRouteSession) getHopsNum() int {
	return tr.hopNum
}

func (tr *TraceRouteSession) runTrace() {
	fmt.Printf("Starting IPv4 tracing (TTL: %d, Timeout: %s)\n", tr.MaxTTL, tr.Timeout)
	tr.traceRoute()
}
