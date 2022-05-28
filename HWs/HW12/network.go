package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Network struct {
	Routers []*Router
}

func NewNetwork(
	configName string,
	routersIPs []string,
) *Network {
	config, err := readConfig(configName)
	if err != nil {
		return nil
	}

	network := make(map[string][]string)
	for routerStr, edges := range config {
		routerInt, err := strconv.Atoi(routerStr)
		if err != nil {
			return nil
		}

		routerEdges := make([]string, len(edges))
		for _, edge := range edges {
			routerEdges = append(routerEdges, routersIPs[edge-1])
		}

		network[routersIPs[routerInt-1]] = routerEdges
	}

	var routers []*Router
	for routerName, edges := range network {
		router := NewRouter(routerName, edges)
		routers = append(routers, router)
	}

	return &Network{Routers: routers}
}

func (n *Network) RIP() {
	step := 0
	isChanged := true

	for isChanged {
		step++
		isChanged = false

		for _, source := range n.Routers {
			for _, destination := range n.Routers {
				if source.Name == destination.Name {
					continue
				}

				for _, edge := range source.Edges {
					change := source.UpdateDistanceTable(
						destination.Name,
						edge,
						destination.Distance(edge),
					)
					if change {
						isChanged = true
					}
				}
			}
			fmt.Printf("Simulation step %v of router %s\n", step, source.Name)
			source.Statistics()
			fmt.Println()
		}
	}
	for _, router := range n.Routers {
		fmt.Printf("Final state of router %s table:\n", router.Name)
		router.Statistics()
		fmt.Println()
	}
}

func readConfig(configName string) (config map[string][]int, err error) {
	jsonFile, err := os.Open(configName)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	return config, nil
}
