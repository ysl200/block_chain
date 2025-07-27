package controller

import (
	"blockchain/global"
	"blockchain/hash"
	"blockchain/node"
	"blockchain/service"
	"encoding/json"
	"net/http"
	"sync"
)

var (
	nodes []*node.Node
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

	no := node.NewNode(id)
	no.CalculateScore(no)

	mu.Lock()
	global.NodesMap[id] = no
	nodes = append(nodes, no)
	hash.AddNode(id)
	mu.Unlock()

	rf := service.NewRaft()
	rf.ElectAnchor(nodes)
}

func (n *NodeController) HandleListNodes(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(global.NodesMap)
}

func (n *NodeController) AddContribution(nodeID string, delta float64) {
	mu.Lock()
	defer mu.Unlock()
	if n, ok := global.NodesMap[nodeID]; ok {
		n.Contribution += delta
		n.CalculateScore(n)
	}
}
