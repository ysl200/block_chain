package service

import (
	"blockchain/global"
	"blockchain/internal/network"
	"sync"
)

var mu sync.Mutex

func GetNodeByID(nodeID string) *network.Node {
	if n, ok := global.NodesMap[nodeID]; ok {
		return n
	}
	return nil
}

func IsCurrentNodeAnchor(nodeID string) bool {
	n := GetNodeByID(nodeID)
	return n != nil && n.IsAnchor
}

// GetAllNodes 获取所有节点
func GetAllNodes() []*network.Node {
	mu.Lock()
	defer mu.Unlock()

	var nodes []*network.Node
	for _, node := range global.NodesMap {
		nodes = append(nodes, node)
	}
	return nodes
}

// AddContribution adds contribution to a node
func AddContribution(nodeID string, delta float64) {
	mu.Lock()
	defer mu.Unlock()
	if n, ok := global.NodesMap[nodeID]; ok {
		n.Contribution += delta
		n.CalculateScore(n)
	}
}
