package handler

import (
	"blockchain/global"
	"blockchain/internal/consensus"
	"blockchain/internal/hash"
	"blockchain/internal/network"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var (
	nodes []*network.Node
	mu    sync.Mutex
)

type NodeController struct {
}

// NewNodeController creates a new NodeController instance
func NewNodeController() *NodeController {
	return &NodeController{}
}

func (n *NodeController) HandleAddNode(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	no := network.NewNode(id)
	no.CalculateScore(no)

	mu.Lock()
	global.NodesMap[id] = no
	nodes = append(nodes, no)
	hash.AddNode(id)
	mu.Unlock()

	rf := consensus.NewRaft()
	anchor := rf.ElectAnchor(nodes)

	// 启动锚节点监听器
	if anchor != nil {
		rf.StartAnchorListener(anchor.ID, global.BlockAssignChan)
	}
}

func (n *NodeController) HandleListNodes(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(global.NodesMap)
}

func (n *NodeController) HandleQueryNode(w http.ResponseWriter, r *http.Request) {
	blockHash := r.URL.Query().Get("block_hash")
	if blockHash == "" {
		http.Error(w, "block_hash is required", http.StatusBadRequest)
		return
	}
	targetNodeID, _ := hash.GetNode(blockHash)
	if targetNodeID == "" {
		http.Error(w, "no node found for the given block hash", http.StatusNotFound)
		return
	}

	rsp := map[string]interface{}{
		"block_hash": blockHash,
		"node_id":    targetNodeID,
	}

	// 5. 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
