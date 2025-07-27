package main

import (
	"blockchain/api/handler"
	"blockchain/global"
	bc "blockchain/internal/blockchain"
	"blockchain/internal/consensus"
	"blockchain/internal/hash"
	"blockchain/internal/network"
	"blockchain/internal/storage"
	"blockchain/web"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.NewSource(time.Now().UnixNano())

	// 创建区块链
	blockchain := bc.NewBlockchain()

	stopTG := make(chan struct{})
	stopM := make(chan struct{})

	// 启动交易生成器和矿工监听器
	blockchain.StartTransactionGenerator(stopTG)
	blockchain.StartMiner(stopM)

	// 模拟运行一段时间
	time.Sleep(3 * time.Second)
	// 停止交易生成器
	close(stopTG)
	log.Println("[交易生成器] 停止交易生成")

	for len(blockchain.PendingTx) >= 3 {
		log.Printf("等待矿工处理剩余交易: %d 笔...\n", len(blockchain.PendingTx))
		time.Sleep(2 * time.Second) // 每隔 2 秒检查一次
	}

	// 停止矿工监听器
	close(stopM)
	log.Println("[矿工监听器] 停止矿工监听")

	// 创建初始节点并启动锚节点监听器
	initializeNodes()

	// 打印区块链信息
	fmt.Println("\n区块链状态:")
	blockchain.PrintBlockchain()

	nodeController := handler.NewNodeController()

	log.Println("服务启动: http://localhost:8080")

	http.HandleFunc("/add", nodeController.HandleAddNode)
	http.HandleFunc("/store", storage.HandleStoreData)
	http.HandleFunc("/list", nodeController.HandleListNodes)
	http.HandleFunc("/query", nodeController.HandleQueryNode)

	web.StartWebServer()

	_ = http.ListenAndServe(":8080", nil)
}

// initializeNodes 初始化节点并启动锚节点监听器
func initializeNodes() {
	// 创建初始节点
	initialNodes := []string{"node1", "node2", "node3", "node4", "node5", "node6", "node7", "node8", "node9", "node10"}

	for _, nodeID := range initialNodes {
		node := network.NewNode(nodeID)
		node.CalculateScore(node)
		log.Printf("[初始化] 创建节点: %s, 节点信息：cpu: %f, memory: %f, disk: %f, bindwitdth: %f\n", node.ID, node.CPU, node.Memory, node.Disk, node.Bandwidth)
		// 添加到全局节点映射
		global.NodesMap[nodeID] = node
		hash.AddNode(nodeID)
	}

	// 创建Raft实例并选举锚节点
	rf := consensus.NewRaft()
	var nodes []*network.Node
	for _, node := range global.NodesMap {
		nodes = append(nodes, node)
	}

	anchor := rf.ElectAnchor(nodes)
	if anchor != nil {
		log.Printf("[初始化] 锚节点 %s 已选举，开始监听区块池", anchor.ID)
		rf.StartAnchorListener(anchor.ID)
	}

	global.AnchorNode = anchor
}
