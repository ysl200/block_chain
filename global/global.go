package global

import "blockchain/node"

var (
	NodesMap = make(map[string]*node.Node)
)
