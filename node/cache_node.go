package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"in-memory-cache-go/cache"
)

type CacheNode struct {
	address        string
	cache          *cache.Cache
	nodes          []string
	lastWriterMap  map[string]string
	syncInProgress bool
	mu             sync.Mutex
}

func NewCacheNode(address string, nodes []string) *CacheNode {
	return &CacheNode{
		address:       address,
		cache:         cache.NewCache(),
		nodes:         nodes,
		lastWriterMap: make(map[string]string),
	}
}

func (node *CacheNode) SetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttlStr := r.URL.Query().Get("ttl")

	if key == "" || value == "" || ttlStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		http.Error(w, "Invalid TTL", http.StatusBadRequest)
		return
	}

	node.cache.Set(key, value, time.Duration(ttl)*time.Second)
	node.mu.Lock()
	node.lastWriterMap[key] = node.address
	node.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cache entry set for key: %s", key)

	node.SyncAfterWrite(key, value, time.Duration(ttl)*time.Second)
}

func (node *CacheNode) SyncAfterWrite(key string, value interface{}, ttl time.Duration) {
	node.mu.Lock()
	node.syncInProgress = true
	node.mu.Unlock()

	for _, nodeAddress := range node.nodes {
		if nodeAddress != node.address {
			go func(addr string) {
				url := fmt.Sprintf("http://%s/cache/broadcast_set?key=%s&value=%v&ttl=%d", addr, key, value, int(ttl.Seconds()))
				_, err := http.Post(url, "application/json", nil)
				if err != nil {
					fmt.Printf("Error syncing to %s: %v\n", addr, err)
					return
				}
			}(nodeAddress)
		}
	}

	node.mu.Lock()
	node.syncInProgress = false
	node.mu.Unlock()
}

func (node *CacheNode) ReceiveSetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttlStr := r.URL.Query().Get("ttl")

	if key == "" || value == "" || ttlStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		http.Error(w, "Invalid TTL", http.StatusBadRequest)
		return
	}

	node.cache.Set(key, value, time.Duration(ttl)*time.Second)
	node.mu.Lock()
	node.lastWriterMap[key] = node.address
	node.mu.Unlock()
	fmt.Printf("Received sync: Cache entry set for key: %s\n", key)
}

func (node *CacheNode) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	node.mu.Lock()
	isSyncing := node.syncInProgress
	lastWriter, isLastWriter := node.lastWriterMap[key]
	node.mu.Unlock()

	if isSyncing && isLastWriter && lastWriter != node.address {
		redirectURL := fmt.Sprintf("http://%s/cache/get?key=%s", lastWriter, key)
		resp, err := http.Get(redirectURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Key not found or expired", http.StatusNotFound)
			return
		}
		io.Copy(w, resp.Body)
	} else {
		value, ok := node.cache.Get(key)
		if !ok {
			http.Error(w, "Key not found or expired", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(value)
	}
}

func (node *CacheNode) Start() {
	http.HandleFunc("/cache/set", node.SetHandler)
	http.HandleFunc("/cache/get", node.GetHandler)
	http.HandleFunc("/cache/broadcast_set", node.ReceiveSetHandler)

	fmt.Printf("Starting cache node at %s\n", node.address)
	http.ListenAndServe(node.address, nil)
}
