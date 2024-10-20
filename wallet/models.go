package wallet

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
