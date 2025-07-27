package config

import "time"

const (
	MinTxToMine     = 3               // 最小交易数，挖矿时需要至少3笔交易
	TxGenInterval   = 1 * time.Second // 交易生成间隔
	MinerCheckDelay = 2 * time.Second // 矿工检查间隔
)
