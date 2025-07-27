package service

import (
	bc "blockchain/block_chain"
	"blockchain/hash"
	"blockchain/node"
	"log"
	"sort"
	"time"
)

type Raft struct{}

// NewRaft New creates a new Raft instance
func NewRaft() *Raft {
	return &Raft{}
}

func (r *Raft) ElectAnchor(nodes []*node.Node) *node.Node {
	if len(nodes) == 0 {
		return nil
	}

	for _, n := range nodes {
		n.CalculateScore(n)
		n.IsAnchor = false // Reset anchor status
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Score > nodes[j].Score // Sort by score in descending order
	})

	// 选出得分最高者作为锚节点
	anchor := nodes[0]
	anchor.IsAnchor = true
	return anchor
}

func (r *Raft) StartAnchorListener(nodeID string) {
	go func() {
		for {
			if !IsCurrentNodeAnchor(nodeID) {
				time.Sleep(30 * time.Second) // 每30秒检查一次
				continue
			}

			blocks := bc.Blockpool
			if len(blocks) == 0 {
				continue
			}
			for _, b := range blocks {
				n, _ := hash.GetNode(b.Hash)
				StoreBlock(n, &b)
				r.AddContribution(nodeID, 10)

				log.Printf("[锚节点] 区块 %d 分发至节点 %s\n", b.Index, nodeID)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (r *Raft) AddContribution(nodeID string, contribution float64) {
	if nodeID == "" || contribution <= 0 {
		return
	}

	n := GetNodeByID(nodeID)
	if n == nil {
		return
	}

	n.Contribution += contribution
}
