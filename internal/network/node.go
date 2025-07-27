package node

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type Node struct {
	ID           string  // Node ID
	CPU          float64 // CPU in percentage (0-100)
	Memory       float64 // Memory in GB
	Disk         float64 // Disk in GB
	Bandwidth    float64 // Bandwidth in Mbps
	Contribution float64 // Contribution score
	Score        float64 // Calculated score
	IsAnchor     bool    // Is this node the anchor?
	Address      string  // Node address (optional, for network communication)
}

// NewNode creates a new Node instance with performance metrics
func NewNode(id string) *Node {
	cpuUsage, memFree, diskFree, bandwidth := getPerformance()

	return &Node{
		ID:           id,
		CPU:          cpuUsage,
		Memory:       memFree,
		Disk:         diskFree,
		Bandwidth:    bandwidth,
		Contribution: 0.0, // Initial contribution is zero
		Score:        0.0, // Initial score is zero
		IsAnchor:     false,
		Address:      "", // Optional address field, can be set later
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

func (n *Node) CalculateScore(node *Node) {
	alpha, beta, gamma, delta := 0.3, 0.3, 0.2, 0.2
	theta, epsilon := 0.6, 0.4

	n.Score = theta*(alpha*node.CPU+beta*node.Disk+gamma*node.Memory+delta*node.Bandwidth) + epsilon*node.Contribution
}

func (n *Node) GetNodeByID(id string) *Node {
	if n.ID == id {
		return n
	}
	return nil
}
