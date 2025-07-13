package service

import (
	"blockchain/global"
	"blockchain/node"
)

func GetNodeByID(nodeID string) *node.Node {
	if n, ok := global.NodesMap[nodeID]; ok {
		return n
	}
	return nil
}

func IsCurrentNodeAnchor(nodeID string) bool {
	n := GetNodeByID(nodeID)
	return n != nil && n.IsAnchor
}
