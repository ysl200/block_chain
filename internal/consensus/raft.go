package consensus

import (
	bc "blockchain/internal/blockchain"
	"blockchain/internal/hash"
	"blockchain/internal/network"
	"blockchain/internal/storage"
	"blockchain/service"
	"log"
	"math/rand"
	"sort"
	"time"
)

type Raft struct{}

// NewRaft New creates a new Raft instance
func NewRaft() *Raft {
	return &Raft{}
}

func (r *Raft) ElectAnchor(nodes []*network.Node) *network.Node {
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

// StartAnchorListener 启动锚节点监听器，监听区块池并分发区块
func (r *Raft) StartAnchorListener(nodeID string, assignChan chan network.BlockAssignInfo) {
	go func() {
		log.Printf("[锚节点监听器] 节点 %s 开始监听区块池...", nodeID)

		lastProcessedIndex := -1 // 记录最后处理的区块索引

		for {
			if !service.IsCurrentNodeAnchor(nodeID) {
				time.Sleep(30 * time.Second) // 每30秒检查一次锚节点状态
				continue
			}

			// 获取所有可用节点
			availableNodes := r.getAvailableNodes()
			if len(availableNodes) == 0 {
				log.Printf("[锚节点] 没有可用的节点进行区块分发")
				time.Sleep(5 * time.Second)
				continue
			}

			// 获取区块池中的新区块
			newBlocks := r.getNewBlocks(lastProcessedIndex)
			if len(newBlocks) == 0 {
				time.Sleep(2 * time.Second)
				continue
			}

			// 分发新区块
			for _, block := range newBlocks {
				r.distributeBlock(block, availableNodes, nodeID, assignChan)
				lastProcessedIndex = block.Index
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

// getAvailableNodes 获取所有可用的节点
func (r *Raft) getAvailableNodes() []*network.Node {
	var availableNodes []*network.Node

	for _, node := range service.GetAllNodes() {
		if node != nil && node.Score > 0 {
			availableNodes = append(availableNodes, node)
		}
	}

	return availableNodes
}

// getNewBlocks 获取新的区块（从上次处理的索引之后）
func (r *Raft) getNewBlocks(lastProcessedIndex int) []bc.Block {
	var newBlocks []bc.Block

	allBlocks := bc.GetAllBlocks()
	for _, block := range allBlocks {
		if block.Index > lastProcessedIndex {
			newBlocks = append(newBlocks, block)
		}
	}

	return newBlocks
}

// distributeBlock 分发区块到不同节点，并通过通道广播分配信息
func (r *Raft) distributeBlock(block bc.Block, availableNodes []*network.Node, anchorNodeID string, assignChan chan network.BlockAssignInfo) {
	// 使用一致性哈希选择目标节点
	targetNodeID, err := hash.GetNode(block.Hash)
	if err != nil {
		// 如果一致性哈希失败，使用负载均衡策略
		targetNodeID = r.selectNodeByLoadBalance(availableNodes)
	}

	// 存储区块到目标节点
	storage.StoreBlock(targetNodeID, &block)

	anchorNode := service.GetNodeByID(anchorNodeID)
	anchorNode.NodeBlockMap[targetNodeID] = append(anchorNode.NodeBlockMap[targetNodeID], block.Hash)

	// 增加锚节点的贡献值
	service.AddContribution(anchorNodeID, 10.0)

	// 增加目标节点的贡献值
	service.AddContribution(targetNodeID, 5.0)

	log.Printf("[锚节点分发] 区块 %d (哈希: %s) 分发至节点 %s",
		block.Index, block.Hash[:8], targetNodeID)

	// 通过通道广播分配信息
	assignChan <- network.BlockAssignInfo{
		Block:        block,
		TargetNodeID: targetNodeID,
	}
}

// selectNodeByLoadBalance 使用负载均衡策略选择节点
func (r *Raft) selectNodeByLoadBalance(nodes []*network.Node) string {
	if len(nodes) == 0 {
		return ""
	}

	// 按分数排序，优先选择分数高的节点
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Score > nodes[j].Score
	})

	// 使用加权随机选择，分数越高的节点被选中的概率越大
	totalScore := 0.0
	for _, node := range nodes {
		totalScore += node.Score
	}

	if totalScore == 0 {
		// 如果所有节点分数都为0，随机选择
		return nodes[rand.Intn(len(nodes))].ID
	}

	// 加权随机选择
	randomValue := rand.Float64() * totalScore
	currentSum := 0.0

	for _, node := range nodes {
		currentSum += node.Score
		if randomValue <= currentSum {
			return node.ID
		}
	}

	// 兜底选择
	return nodes[0].ID
}

func (r *Raft) AddContribution(nodeID string, contribution float64) {
	if nodeID == "" || contribution <= 0 {
		return
	}

	service.AddContribution(nodeID, contribution)
}
