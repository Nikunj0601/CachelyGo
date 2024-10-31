package main

import (
	"flag"
	"fmt"

	"in-memory-cache-go/loadbalancer"
	"in-memory-cache-go/node"
	"strings"
)

func splitNodes(nodes string) []string {
	return strings.Split(nodes, ",")
}

func main() {
	mode := flag.String("mode", "node", "Start mode: 'node' or 'load_balancer'")
	port := flag.String("port", "8081", "Port to run the server on")
	nodes := flag.String("nodes", "localhost:8081,localhost:8082,localhost:8083", "Comma-separated list of cache nodes")
	flag.Parse()

	nodeList := splitNodes(*nodes)

	fmt.Println(*mode, *port, nodeList)

	if *mode == "load_balancer" {
		fmt.Println("Starting load balancer on port 8080")
		loadBalancer := loadbalancer.NewLoadBalancer(nodeList)
		loadBalancer.Start()
	} else {
		fmt.Printf("Starting cache node on port %s\n", *port)
		cacheNode := node.NewCacheNode(fmt.Sprintf("localhost:%s", *port), nodeList)
		cacheNode.Start()
	}
}
