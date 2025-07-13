package contribution

import "blockchain/service"

func AddContribution(nodeID string, contribution float64) {
	// 获取节点
	node := service.GetNodeByID(nodeID)
	if node == nil {
		return // 节点不存在
	}
	node.Contribution += contribution
}
