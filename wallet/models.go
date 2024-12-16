package wallet

import "tajfi-server/wallet/tapd"

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

// Order represents a sell order.
type Order struct {
	AssetID             string        `json:"asset_id"`
	AmountToSell        int64         `json:"amount_to_sell"`
	AmountSatsToReceive int64         `json:"amount_sats_to_receive"`
	Outpoint            tapd.Outpoint `json:"outpoint"`
	//VirtualOutpoint     tapd.Outpoint `json:"virtual_outpoint"`
	VirtualPSBT       string   `json:"virtual_psbt"`
	AnchorPSBT        string   `json:"anchor_psbt"`
	PassiveAssetPSBTs []string `json:"passive_asset_psbts"`
}
