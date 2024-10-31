package loadbalancer

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	nodes []string
	index int
	mu    sync.Mutex
}

func NewLoadBalancer(nodes []string) *LoadBalancer {
	return &LoadBalancer{nodes: nodes}
}

func (lb *LoadBalancer) GetNode() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	node := lb.nodes[lb.index]
	lb.index = (lb.index + 1) % len(lb.nodes)
	return node
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetNode := lb.GetNode()
	url := fmt.Sprintf("http://%s%s?%s", targetNode, r.URL.Path, r.URL.RawQuery)
	proxyReq, _ := http.NewRequest(r.Method, url, r.Body)
	proxyReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (lb *LoadBalancer) Start() {
	fmt.Println("Load balancer started on localhost:8080")
	http.ListenAndServe(":8080", lb)
}
