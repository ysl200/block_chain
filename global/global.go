package global

import (
	block_chain "blockchain/internal/blockchain"
	"blockchain/internal/network"
)

var (
	NodesMap        = make(map[string]*network.Node) // 存储节点信息的全局映射
	Blockpool       = make([]block_chain.Block, 0)   // 全局区块池，用于存储待处理的区块
	AnchorNode      *network.Node                    // 当前锚节点
	BlockAssignChan chan network.BlockAssignInfo
)
