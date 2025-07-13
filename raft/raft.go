package raft

import (
	"blockchain/node"
	"blockchain/service"
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
			if !service.IsCurrentNodeAnchor(nodeID) {
				time.Sleep(30 * time.Second) // 每30秒检查一次
				continue
			}

			//block := blockpoll.GetNewBlock()
			//if block == nil {
			//	target := hash.GetNode()
			//	storage.StoreBlock(target, block)
			//
			//}
		}
	}()
}
