package main

import (
	"fmt"
	"math"
)

type Router struct {
	Name          string
	Edges         []string
	DistanceTable map[string]*NextHopAndMetric
}

type NextHopAndMetric struct {
	NextHop string
	Metric  float64
}

func NewRouter(
	name string,
	edges []string,
) *Router {
	distanceTable := make(map[string]*NextHopAndMetric)

	for _, edge := range edges {
		if edge == "" {
			continue
		}
		distanceTable[edge] = &NextHopAndMetric{
			NextHop: edge,
			Metric:  1,
		}
	}

	return &Router{
		Name:          name,
		Edges:         edges,
		DistanceTable: distanceTable,
	}
}

func (r *Router) Distance(routerName string) float64 {
	distance, ok := r.DistanceTable[routerName]
	if !ok {
		return math.Inf(1)
	}

	return distance.Metric
}

func (r *Router) UpdateDistanceTable(
	routerName string,
	nextHop string,
	distance float64,
) bool {
	if !math.IsInf(distance, 1) {
		currDistance, ok := r.DistanceTable[routerName]
		if !ok || currDistance.Metric > distance+1 {
			r.DistanceTable[routerName] = &NextHopAndMetric{
				NextHop: nextHop,
				Metric:  distance + 1,
			}
			return true
		}
	}

	return false
}

func (r *Router) Statistics() {
	fmt.Printf("%-25s%-25s%-25s%-25s\n",
		"[Source IP]", "[Destination IP]", "[Next Hop]", "[Metric]")

	for routerName, nextHopAndMetric := range r.DistanceTable {
		fmt.Printf("%-25s%-25s%-25s%-25v\n",
			r.Name, routerName, nextHopAndMetric.NextHop, nextHopAndMetric.Metric)
	}
}
