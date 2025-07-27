package handler

import (
	"blockchain/global"
	"blockchain/internal/consensus"
	"blockchain/internal/hash"
	"blockchain/internal/network"
	"encoding/json"
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
		rf.StartAnchorListener(anchor.ID)
	}
}

func (n *NodeController) HandleListNodes(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(global.NodesMap)
}
