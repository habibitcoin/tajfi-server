package tapd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"tajfi-server/wallet/lnd"
)

// DecodeAddrResponse represents the full decoded asset response from Tapd.
type DecodeAddrResponse struct {
	Encoded          string `json:"encoded"`
	AssetID          string `json:"asset_id"`
	AssetType        string `json:"asset_type"`
	Amount           string `json:"amount"`
	GroupKey         string `json:"group_key"`
	ScriptKey        string `json:"script_key"`
	InternalKey      string `json:"internal_key"`
	TapscriptSibling string `json:"tapscript_sibling"`
	TaprootOutputKey string `json:"taproot_output_key"`
	ProofCourierAddr string `json:"proof_courier_addr"`
	AssetVersion     string `json:"asset_version"`
	AddressVersion   string `json:"address_version"`
}

// DecodeAddr decodes a Taproot Asset address into the full asset response.
func DecodeAddr(tapdHost, macaroon, address string) (*DecodeAddrResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/addrs/decode", tapdHost)

	// Prepare the request payload.
	payload := map[string]string{"addr": address}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Disable TLS verification (for testing).
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Create the HTTP request.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tapd RPC error: %s", resp.Status)
	}

	// Parse the response into DecodeAddrResponse.
	var decodedAddr DecodeAddrResponse
	if err := json.NewDecoder(resp.Body).Decode(&decodedAddr); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &decodedAddr, nil
}

type NewAddressPayload struct {
	AssetID     string                  `json:"asset_id"`
	Amt         int                     `json:"amt"`
	ScriptKey   map[string]interface{}  `json:"script_key"`
	InternalKey lnd.InternalKeyResponse `json:"internal_key"`
}

// CallNewAddress sends the payload to the Tapd NewAddress RPC.
func CallNewAddress(tapdHost, macaroon string, payload NewAddressPayload) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/addrs", tapdHost)

	// Disable TLS verification (use with caution)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tapd RPC error: %s", resp.Status)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FundVirtualPSBT sends a request to Tapd to fund a virtual PSBT using the invoice.
func FundVirtualPSBT(tapdHost, macaroon, invoice string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/fund", tapdHost)

	// Prepare the payload
	requestBody := map[string]interface{}{
		"raw": map[string]interface{}{
			"recipients": map[string]int{
				invoice: 0, // Use the invoice as the recipient
			},
		},
	}
	payloadBytes, _ := json.Marshal(requestBody)

	// Disable TLS verification (for testing purposes)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tapd RPC error: %s", resp.Status)
	}

	// Parse the response
	var fundedPsbt map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&fundedPsbt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return fundedPsbt, nil
}

// AnchorVirtualPSBTParams holds the parameters for the anchor request.
type AnchorVirtualPSBTParams struct {
	VirtualPSBTs []string // Raw bytes for virtual PSBTs
	TapdHost     string
	Macaroon     string
}

// AssetTransferResponse represents the response from the Tapd API.
type AssetTransferResponse struct {
	TransferTimestamp  string                   `json:"transfer_timestamp"`
	AnchorTxHash       string                   `json:"anchor_tx_hash"`
	AnchorTxHeightHint int                      `json:"anchor_tx_height_hint"`
	AnchorTxChainFees  string                   `json:"anchor_tx_chain_fees"`
	Inputs             []map[string]interface{} `json:"inputs"`
	Outputs            []map[string]interface{} `json:"outputs"`
	AnchorTxBlockHash  map[string]interface{}   `json:"anchor_tx_block_hash"`
}

func AnchorVirtualPSBT(params AnchorVirtualPSBTParams) (*AssetTransferResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/anchor", params.TapdHost)

	// Prepare the request payload
	payload := map[string]interface{}{
		"virtual_psbts": params.VirtualPSBTs,
	}
	payloadBytes, _ := json.Marshal(payload)

	// Disable TLS verification (for testing)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", params.Macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tapd RPC error: %s", resp.Status)
	}

	// Parse the response into the AssetTransferResponse struct
	var assetTransfer AssetTransferResponse
	err = json.NewDecoder(resp.Body).Decode(&assetTransfer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &assetTransfer, nil
}
