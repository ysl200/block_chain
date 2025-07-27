package network

// HealthStatus 定义节点健康状态
type HealthStatus struct {
	Status      string  // 状态描述 (健康/警告/危险)
	Score       float64 // 健康评分 (0-100)
	CPUUsage    float64 // CPU使用率
	MemoryUsage float64 // 内存使用率
	DiskUsage   float64 // 磁盘使用率
	NetActivity float64 // 网络活动度
}

// CheckHealth 检查节点健康状态
func (n *Node) CheckHealth() HealthStatus {
	// 获取当前性能指标
	cpuUsage, memFree, diskFree, bandwidth := getPerformance()

	// 计算使用率百分比
	cpuUsagePct := 100 - cpuUsage
	memUsagePct := (1 - memFree/(memFree+n.Memory)) * 100
	diskUsagePct := (1 - diskFree/(diskFree+n.Disk)) * 100

	// 计算健康评分 (权重可调整)
	healthScore := 0.3*(100-cpuUsagePct) + 0.2*(100-memUsagePct) +
		0.3*(100-diskUsagePct) + 0.2*(bandwidth/1000) // 假设1Gbps为满分

	// 确定状态级别
	status := "健康"
	if healthScore < 60 {
		status = "危险"
	} else if healthScore < 80 {
		status = "警告"
	}

	return HealthStatus{
		Status:      status,
		Score:       healthScore,
		CPUUsage:    cpuUsagePct,
		MemoryUsage: memUsagePct,
		DiskUsage:   diskUsagePct,
		NetActivity: bandwidth,
	}
}

// IsHealthy 判断节点是否健康(简化版)
func (n *Node) IsHealthy() bool {
	health := n.CheckHealth()
	return health.Score >= 60 // 评分≥60视为健康
}
