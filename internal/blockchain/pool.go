package block_chain

import (
	"sync"
)

var (
	mu sync.Mutex
)

func AddBlock(block Block) {
	mu.Lock()
	defer mu.Unlock()
	Blocks = append(Blocks, block)
}

func GetBlock() *Block {
	mu.Lock()
	defer mu.Unlock()
	if len(Blocks) == 0 {
		return nil
	}
	block := Blocks[0]
	Blocks = Blocks[1:]
	return &block
}

func GetAllBlocks() []Block {
	mu.Lock()
	defer mu.Unlock()
	return append([]Block{}, Blocks...) // 返回一个副本，避免外部修改
}
