package hash

import (
	"github.com/stathat/consistent"
	"sync"
)

var (
	ring = consistent.New()
	mu   sync.Mutex
)

func AddNode(nodeID string) {
	mu.Lock()
	defer mu.Unlock()

	// Add the node to the consistent hash ring
	ring.Add(nodeID)
}

func RemoveNode(nodeID string) {
	mu.Lock()
	defer mu.Unlock()

	// Remove the node from the consistent hash ring
	ring.Remove(nodeID)
}

func GetNode(key string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	// Get the node responsible for the given key
	nodeID, err := ring.Get(key)
	if err != nil {
		return "", err // Handle error appropriately
	}

	return nodeID, nil
}
