package tapd

import "fmt"

// Structs for balances, transfers, and UTXOs
type AssetGenesis struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
}

type AssetBalance struct {
	AssetGenesis       AssetGenesis `json:"asset_genesis"`
	Balance            string       `json:"balance"`
	UnconfirmedBalance string       `json:"unconfirmed_balance,omitempty"`
}

type WalletBalancesResponse struct {
	AssetBalances map[string]AssetBalance `json:"asset_balances"`
}

type Asset struct {
	AssetGenesis AssetGenesis `json:"asset_genesis"`
	Amount       string       `json:"amount"`
	ScriptKey    string       `json:"script_key"`
}

type ManagedUtxo struct {
	Outpoint string  `json:"out_point"`
	Assets   []Asset `json:"assets"`
}

type GetUtxosResponse struct {
	ManagedUtxos map[string]ManagedUtxo `json:"managed_utxos"`
}

type AnchorTxBlockHash struct {
	Hash    string `json:"hash"`
	HashStr string `json:"hash_str"`
}

type AssetTransferResponse struct {
	Inputs             []TransferInput   `json:"inputs"`
	Outputs            []TransferOutput  `json:"outputs"`
	AnchorTxHash       string            `json:"anchor_tx_hash"`
	AnchorTxHeightHint int               `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  string            `json:"anchor_tx_chain_fees"`
	AnchorTxBlockHash  AnchorTxBlockHash `json:"anchor_tx_block_hash"`
	TransferTimestamp  string            `json:"transfer_timestamp"`
}

type TransferInput struct {
	AnchorPoint string `json:"anchor_point"`
	AssetID     string `json:"asset_id"`
	ScriptKey   string `json:"script_key"`
	Amount      string `json:"amount"`
}

type TransferOutput struct {
	Anchor    Anchor `json:"anchor"`
	ScriptKey string `json:"script_key"`
	Amount    string `json:"amount"`
}

type Anchor struct {
	Outpoint         string `json:"outpoint"`
	Value            string `json:"value"`
	InternalKey      string `json:"internal_key"`
	TaprootAssetRoot string `json:"taproot_asset_root"`
}

type AssetTransfersResponse struct {
	Transfers []AssetTransferResponse `json:"transfers"`
}

// GetBalances fetches wallet balances using a GET request.
func (c *tapdClient) GetBalances(tapdHost, macaroon string) (*WalletBalancesResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/assets/balance?asset_id=true&include_leased=true", tapdHost)

	var balances WalletBalancesResponse
	err := c.makeGetRequest(url, macaroon, &balances)
	return &balances, err
}

// GetUtxos fetches wallet UTXOs using a GET request.
func (c *tapdClient) GetUtxos(tapdHost, macaroon string) (*GetUtxosResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/assets/utxos", tapdHost)

	var utxos GetUtxosResponse
	err := c.makeGetRequest(url, macaroon, &utxos)
	return &utxos, err
}

// GetTransfers fetches asset transfers using a GET request.
func (c *tapdClient) GetTransfers(tapdHost, macaroon string) (AssetTransfersResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/assets/transfers", tapdHost)

	var transfers AssetTransfersResponse
	err := c.makeGetRequest(url, macaroon, &transfers)
	return transfers, err
}
