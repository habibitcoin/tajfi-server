package tapd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AssetTransferResponse represents the response from the Tapd API.
type TransferInput struct {
	AnchorPoint string `json:"anchor_point"`
	AssetID     string `json:"asset_id"`
	ScriptKey   string `json:"script_key"`
	Amount      string `json:"amount"`
}

// Anchor represents the nested anchor object within TransferOutput.
type Anchor struct {
	Outpoint         string `json:"outpoint"`
	Value            string `json:"value"`
	InternalKey      string `json:"internal_key"`
	TaprootAssetRoot string `json:"taproot_asset_root"`
	MerkleRoot       string `json:"merkle_root"`
	TapscriptSibling string `json:"tapscript_sibling"`
	NumPassiveAssets int    `json:"num_passive_assets"`
}

type TransferOutput struct {
	Anchor              Anchor `json:"anchor"`
	ScriptKey           string `json:"script_key"`
	ScriptKeyIsLocal    bool   `json:"script_key_is_local"`
	Amount              string `json:"amount"`
	NewProofBlob        string `json:"new_proof_blob"`
	SplitCommitRootHash string `json:"split_commit_root_hash"`
	OutputType          string `json:"output_type"`
	AssetVersion        string `json:"asset_version"`
	LockTime            string `json:"lock_time"`
	RelativeLockTime    string `json:"relative_lock_time"`
	ProofDeliveryStatus string `json:"proof_delivery_status,omitempty"`
}

type AnchorTxBlockHash struct {
	Hash    string `json:"hash"`
	HashStr string `json:"hash_str"`
}

type AssetTransferResponse struct {
	TransferTimestamp  string            `json:"transfer_timestamp"`
	AnchorTxHash       string            `json:"anchor_tx_hash"`
	AnchorTxHeightHint int               `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  string            `json:"anchor_tx_chain_fees"`
	AnchorTxBlockHash  AnchorTxBlockHash `json:"anchor_tx_block_hash"`
	Inputs             []TransferInput   `json:"inputs"`
	Outputs            []TransferOutput  `json:"outputs"`
}

type AssetGenesis struct {
	GenesisPoint string `json:"genesis_point"`
	Name         string `json:"name"`
	MetaHash     string `json:"meta_hash"`
	AssetID      string `json:"asset_id"`
	AssetType    string `json:"asset_type"`
	OutputIndex  int    `json:"output_index"`
}

// AssetBalance represents an individual asset's balance.
type AssetBalance struct {
	AssetGenesis       AssetGenesis `json:"asset_genesis"`
	Balance            string       `json:"balance"`
	UnconfirmedBalance string       `json:"unconfirmed_balance,omitempty"`
}

// WalletBalancesResponse represents the response structure for wallet balances.
type WalletBalancesResponse struct {
	AssetBalances map[string]AssetBalance `json:"asset_balances"`
}

// GetBalances interacts with the tapd daemon to retrieve wallet balances.
func (c *tapdClient) GetBalances(tapdHost, macaroon string) (*WalletBalancesResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/assets/balance?asset_id=true&include_leased=true", tapdHost)
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tapd daemon error: %s", resp.Status)
	}

	// Decode the response body into WalletBalancesResponse
	var balances WalletBalancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &balances, nil
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

// WalletBalancesResponse represents the response structure for wallet balances.
type GetUtxosResponse struct {
	ManagedUtxos map[string]ManagedUtxo `json:"managed_utxos"`
}

// GetBalances interacts with the tapd daemon to retrieve wallet UTXOs.
func (c *tapdClient) GetUtxos(tapdHost, macaroon string) (*GetUtxosResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/assets/utxos", tapdHost)
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tapd daemon error: %s", resp.Status)
	}

	// Decode the response body into WalletBalancesResponse
	var utxos GetUtxosResponse
	if err := json.NewDecoder(resp.Body).Decode(&utxos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &utxos, nil
}

type AssetTransfersResponse struct {
	Transfers []AssetTransferResponse `json:"transfers"`
}

func (c *tapdClient) GetTransfers(tapdHost, macaroon string) (transfers AssetTransfersResponse, err error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/assets/transfers", tapdHost)
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return transfers, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return transfers, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return transfers, fmt.Errorf("tapd daemon error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&transfers); err != nil {
		return transfers, fmt.Errorf("failed to decode response: %v", err)
	}

	return transfers, nil
}
