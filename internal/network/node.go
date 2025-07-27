package network

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type Node struct {
	ID           string
	CPU          float64
	Memory       float64
	Disk         float64
	Bandwidth    float64
	Contribution float64
	Score        float64
	IsAnchor     bool
	Address      string
	NodeBlockMap map[string][]string
	LastHealth   HealthStatus // 新增健康状态记录
}

// 更新NewNode函数
func NewNode(id string) *Node {
	cpuUsage, memFree, diskFree, bandwidth := getPerformance()

	node := &Node{
		ID:           id,
		CPU:          cpuUsage,
		Memory:       memFree,
		Disk:         diskFree,
		Bandwidth:    bandwidth,
		Contribution: 0.0,
		Score:        0.0,
		IsAnchor:     false,
		Address:      "",
		NodeBlockMap: make(map[string][]string),
	}

	// 初始化健康状态
	node.LastHealth = node.CheckHealth()

	return node
}

// CalculateScore 更新CalculateScore时考虑健康状态
func (n *Node) CalculateScore(node *Node) {
	alpha, beta, gamma, delta := 0.3, 0.3, 0.2, 0.2
	theta, epsilon := 0.6, 0.4

	// 获取当前健康评分(0-100)并归一化到0-1
	healthScore := node.CheckHealth().Score / 100

	n.Score = theta*(alpha*node.CPU+beta*node.Disk+gamma*node.Memory+delta*node.Bandwidth) +
		epsilon*(node.Contribution*healthScore)
}

// StoreBlock 模拟存储区块并扣减磁盘空间
func (n *Node) StoreBlock(blockID string, blockSizeGB float64) error {
	// 检查磁盘空间是否足够
	if blockSizeGB > n.Disk {
		return fmt.Errorf("磁盘空间不足，需要 %.2fGB，可用 %.2fGB", blockSizeGB, n.Disk)
	}

	// 扣减磁盘空间
	n.Disk -= blockSizeGB

	// 更新区块映射
	if _, exists := n.NodeBlockMap[blockID]; !exists {
		n.NodeBlockMap[blockID] = []string{}
	}

	return nil
}

// RemoveBlock 模拟删除区块并释放磁盘空间
func (n *Node) RemoveBlock(blockID string, blockSizeGB float64) {
	if _, exists := n.NodeBlockMap[blockID]; exists {
		delete(n.NodeBlockMap, blockID)
		n.Disk += blockSizeGB
	}
}

func getPerformance() (float64, float64, float64, float64) {
	cpuPercent, _ := cpu.Percent(0, false)
	memStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")
	netStat, _ := net.IOCounters(false)

	cpuUsage := 100 - cpuPercent[0]             // 预留越多性能越高
	memFree := float64(memStat.Available) / 1e9 // GB
	diskFree := float64(diskStat.Free) / 1e9
	bandwidth := float64(netStat[0].BytesRecv+netStat[0].BytesSent) / 1e6

	return cpuUsage, memFree, diskFree, bandwidth
}
