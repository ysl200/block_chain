package block_chain

type Transaction struct {
	ID          string  `json:"id"`
	Sender      string  `json:"sender"`
	Recipient   string  `json:"recipient"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
	Description string  `json:"description"`
}
