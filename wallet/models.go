package wallet

type Transfer struct {
	Txid         string `json:"txid"`
	Timestamp    string `json:"timestamp"`
	Height       int    `json:"height"`
	AssetID      string `json:"asset_id"`
	Type         string `json:"type"`
	Amount       uint64 `json:"amount"`
	Status       string `json:"status"`        // confirmed or unconfirmed
	ChangeAmount uint64 `json:"change_amount"` // only relevant for sends
}

type Wallet struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

// SendStartParams holds parameters needed to start a send operation.
type ReceiveParams struct {
	PubKey       string
	AssetID      string
	Amount       int
	LNDHost      string
	LNMacaroon   string
	TapdHost     string
	TapdMacaroon string
}
