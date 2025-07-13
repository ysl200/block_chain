package storage

import (
	"blockchain/controller"
	"blockchain/hash"
	"fmt"
	"net/http"
)

var (
	data = make(map[string][]block.Block)
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

func StoreBlock(nodeID string, block *blockchain.Block)
