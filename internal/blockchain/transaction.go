package block_chain

import (
	"math/rand"
	"strconv"
	"time"
)

type Transaction struct {
	ID          string  `json:"id"`
	Sender      string  `json:"sender"`
	Recipient   string  `json:"recipient"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
	Description string  `json:"description"`
}

var (
	pool []*Transaction
)

func AddTransaction(tx Transaction) {
	mu.Lock()
	defer mu.Unlock()
	pool = append(pool, &tx)
}

func GetTransactions(n int) []*Transaction {
	mu.Lock()
	defer mu.Unlock()

	if len(pool) < n {
		return nil
	}

	txs := pool[:n]
	pool = pool[n:]
	return txs
}

func StartGenerating() {
	go func() {
		for {
			// 模拟生成交易
			tx := Transaction{
				ID:          "tx" + strconv.Itoa(rand.Intn(10000)),
				Sender:      "user" + strconv.Itoa(rand.Intn(10)),
				Recipient:   "user" + strconv.Itoa(rand.Intn(10)),
				Amount:      float64(rand.Intn(1000) + 1),
				Timestamp:   time.Now().Unix(),
				Description: "Test transaction",
			}
			AddTransaction(tx)
			time.Sleep(100 * time.Millisecond) // 每100ms生成一个交易
		}
	}()
}
