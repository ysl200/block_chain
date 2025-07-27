package service

import (
	"blockchain/block_chain"
	"blockchain/controller"
	"blockchain/hash"
	"fmt"
	"net/http"
	"sync"
)

var (
	data = make(map[string][]block_chain.Block)
	mu   sync.Mutex
)

func HandleStoreData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	nodeID, err := hash.GetNode(key)
	if err != nil {
		http.Error(w, "failed to get node: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nodeController := controller.NodeController{}
	nodeController.AddContribution(nodeID, 1.0)
	w.Write([]byte(fmt.Sprintf("Key %s stored on node %s\n", key, nodeID)))
}

func StoreBlock(nodeID string, block *block_chain.Block) {
	mu.Lock()
	defer mu.Unlock()
	data[nodeID] = append(data[nodeID], *block)
}

func GetNodeBlocks(nodeID string) ([]block_chain.Block, error) {
	mu.Lock()
	defer mu.Unlock()
	blocks, exists := data[nodeID]
	if !exists {
		return nil, fmt.Errorf("no blocks found for node %s", nodeID)
	}
	return blocks, nil
}

func GetAllData() map[string][]block_chain.Block {
	mu.Lock()
	defer mu.Unlock()
	// 返回一个副本，避免外部修改
	dataCopy := make(map[string][]block_chain.Block)
	for k, v := range data {
		dataCopy[k] = append([]block_chain.Block{}, v...)
	}
	return dataCopy
}
