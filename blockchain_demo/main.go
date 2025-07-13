package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MinTxToMine     = 3               // 最小交易数，挖矿时需要至少3笔交易
	TxGenInterval   = 2 * time.Second // 交易生成间隔
	MinerCheckDelay = 1 * time.Second // 矿工检查间隔
)

type PowResult struct {
	Block Block  `json:"block"`
	Hash  string `json:"hash"`
}

// Transaction 定义交易结构
type Transaction struct {
	ID          string  `json:"id"`
	Sender      string  `json:"sender"`
	Recipient   string  `json:"recipient"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
	Description string  `json:"description"`
}

// Block 定义区块结构
type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prevHash"`
	Hash         string        `json:"hash"`
	Nonce        int           `json:"nonce"`
	Miner        string        `json:"miner"`
}

// Blockchain 定义区块链结构
type Blockchain struct {
	Chain         []Block        `json:"chain"`
	PendingTx     []Transaction  `json:"pendingTx"`
	Difficulty    int            `json:"difficulty"`
	Reward        float64        `json:"reward"`
	Addresses     []string       `json:"addresses"`
	Products      []string       `json:"products"`
	Contributions map[string]int `json:"contributions"`
}

// NewBlockchain 创建新的区块链
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Difficulty: 4, // 设置为4，哈希只需前4位为0
		Reward:     10.0,
		Addresses: []string{
			"Alice", "Bob", "Charlie", "David", "Eve",
			"Frank", "Grace", "Heidi", "Ivan", "Judy",
		},
		Products: []string{
			"Laptop", "Phone", "Tablet", "Camera", "Headphones",
			"Monitor", "Keyboard", "Mouse", "Printer", "Router",
		},
		Contributions: make(map[string]int),
	}

	for _, addr := range bc.Addresses {
		bc.Contributions[addr] = 1 // 初始化每个地址的贡献值为1
	}

	bc.CreateGenesisBlock()
	return bc
}

// CreateGenesisBlock 创建创世区块
func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Nonce:        0,
		Miner:        "Genesis",
	}
	genesisBlock.Hash = bc.CalculateHash(genesisBlock)
	bc.Chain = append(bc.Chain, genesisBlock)
}

// GenerateRandomTransaction 生成随机交易
func (bc *Blockchain) GenerateRandomTransaction() Transaction {
	rand.NewSource(time.Now().UnixNano())
	sender := bc.Addresses[rand.Intn(len(bc.Addresses))]
	recipient := bc.Addresses[rand.Intn(len(bc.Addresses))]
	for recipient == sender {
		recipient = bc.Addresses[rand.Intn(len(bc.Addresses))]
	}

	return Transaction{
		ID:          hex.EncodeToString([]byte(strconv.Itoa(rand.Int()))),
		Sender:      sender,
		Recipient:   recipient,
		Amount:      float64(rand.Intn(1000)+1) + rand.Float64(),
		Timestamp:   time.Now().Unix(),
		Description: fmt.Sprintf("Purchase of %s", bc.Products[rand.Intn(len(bc.Products))]),
	}
}

// AddTransaction 添加交易到待处理交易池
func (bc *Blockchain) AddTransaction(tx Transaction) {
	bc.PendingTx = append(bc.PendingTx, tx)
}

// CalculateHash 优化后的哈希计算函数
func (bc *Blockchain) CalculateHash(block Block) string {
	// 只序列化影响哈希的关键字段，提高效率
	type tempBlock struct {
		Index        int
		Timestamp    int64
		Transactions []Transaction
		PrevHash     string
		Nonce        int
	}

	tb := tempBlock{
		Index:        block.Index,
		Timestamp:    block.Timestamp,
		Transactions: block.Transactions,
		PrevHash:     block.PrevHash,
		Nonce:        block.Nonce,
	}

	blockBytes, _ := json.Marshal(tb)
	h := sha256.Sum256(blockBytes)
	return hex.EncodeToString(h[:])
}

// GetRandomMinerByContribution 权重随机矿工
func (bc *Blockchain) GetRandomMinerByContribution() string {
	total := 0
	for _, v := range bc.Contributions {
		total += v
	}
	r := rand.Intn(total)
	acc := 0
	for miner, c := range bc.Contributions {
		acc += c
		if r < acc {
			return miner
		}
	}
	return bc.Addresses[rand.Intn(len(bc.Addresses))] // 默认返回随机地址
}

// ProofOfWorkParallel 多线程挖矿
func (bc *Blockchain) ProofOfWorkParallel(lastBlock Block, workerCount int) (Block, string) {

	resultChan := make(chan PowResult, 1)
	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			block := Block{
				Index:        lastBlock.Index + 1,
				Timestamp:    time.Now().Unix(),
				Transactions: append([]Transaction(nil), bc.PendingTx...),
				PrevHash:     lastBlock.Hash,
				Nonce:        rand.Intn(10000), // 随机初始Nonce值
				Miner:        bc.GetRandomMinerByContribution(),
			}

			rewardTx := Transaction{
				ID:          hex.EncodeToString([]byte(strconv.Itoa(rand.Int()))),
				Sender:      "System",
				Recipient:   block.Miner,
				Amount:      bc.Reward,
				Timestamp:   time.Now().Unix(),
				Description: "Mining reward",
			}
			block.Transactions = append(block.Transactions, rewardTx)

			for {
				select {
				case <-stopChan:
					return
				default:
					hash := bc.CalculateHash(block)
					if strings.HasPrefix(hash, strings.Repeat("0", bc.Difficulty)) {
						select {
						case resultChan <- PowResult{block, hash}:
							close(stopChan) // 发送结果并关闭停止信号
						default:
						}
						return
					}
					block.Nonce++
				}
			}
		}()
	}

	wg.Wait()
	select {
	case result := <-resultChan:
		return result.Block, result.Hash
	default:
		fmt.Println("挖矿超时，未找到有效哈希")
		return Block{}, ""
	}
}

// MineBlock 挖矿生成新区块
func (bc *Blockchain) MineBlock() Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	newBlock, hash := bc.ProofOfWorkParallel(lastBlock, 4) // 使用4个工作线程进行挖矿
	if hash == "" {
		fmt.Println("挖矿失败")
		return Block{}
	}

	newBlock.Hash = hash
	bc.Chain = append(bc.Chain, newBlock)
	bc.PendingTx = []Transaction{}
	bc.Contributions[newBlock.Miner]++

	return newBlock
}

// PrintBlockchain 打印区块链信息
func (bc *Blockchain) PrintBlockchain() {
	for _, block := range bc.Chain {
		fmt.Printf("\n区块 %d:\n", block.Index)
		fmt.Printf("  时间戳: %d\n", block.Timestamp)
		fmt.Printf("  哈希: %s\n", block.Hash)
		fmt.Printf("  前哈希: %s\n", block.PrevHash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Printf("  矿工: %s\n", block.Miner)
		fmt.Printf("  交易数: %d\n", len(block.Transactions))
		for _, tx := range block.Transactions {
			fmt.Printf("    %s -> %s: %.2f (%s)\n", tx.Sender, tx.Recipient, tx.Amount, tx.Description)
		}
	}
}

func StartTransactionGenerator(bc *Blockchain, stop chan struct{}) {
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				tx := bc.GenerateRandomTransaction()
				bc.AddTransaction(tx)
				fmt.Printf("[交易生成] %s -> %s %.2f (%s)\n",
					tx.Sender, tx.Recipient, tx.Amount, tx.Description)
				time.Sleep(TxGenInterval)
			}
		}
	}()
}

func StartMiner(bc *Blockchain, stop <-chan struct{}) {
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				if len(bc.PendingTx) >= MinTxToMine {
					fmt.Printf("\n[矿工检测] 当前交易池中有 %d 笔交易，开始挖矿...\n", len(bc.PendingTx))
					start := time.Now()
					block := bc.MineBlock()
					if block.Index > 0 {
						fmt.Printf("[挖矿完成] 区块 %d 被 %s 挖出，耗时 %v\n\n", block.Index, block.Miner, time.Since(start))
					}
				}
				time.Sleep(MinerCheckDelay)
			}
		}
	}()
}

func main() {
	rand.NewSource(time.Now().UnixNano())

	// 创建区块链
	blockchain := NewBlockchain()

	stop := make(chan struct{})

	// 启动交易生成器和矿工监听器
	StartTransactionGenerator(blockchain, stop)
	StartMiner(blockchain, stop)

	// 模拟运行一段时间
	time.Sleep(30 * time.Second)
	// 停止交易生成器和矿工
	close(stop)

	// 打印区块链信息
	fmt.Println("\n区块链状态:")
	blockchain.PrintBlockchain()
}
