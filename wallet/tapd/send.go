package tapd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type SendAssetsResponse struct {
	FundedPSBT        string   `json:"funded_psbt"`
	ChangeOutputIndex int      `json:"change_output_index"`
	PassiveAssetPsbts []string `json:"passive_asset_psbts"`
	SighashHexToSign  string   `json:"sighash_hex_to_sign"`
}

type FundVirtualPSBTResponse struct {
	FundedPSBT        string   `json:"funded_psbt"`
	ChangeOutputIndex int      `json:"change_output_index"`
	PassiveAssetPsbts []string `json:"passive_asset_psbts"`
	SighashHexToSign  string   `json:"sighash_hex_to_sign"`
}

type PrevIds struct {
	Inputs []PrevId `json:"inputs,omitempty"`
}

type Outpoint struct {
	Txid        string `json:"txid"`
	OutputIndex int    `json:"output_index"`
}

type PrevId struct {
	Outpoint  Outpoint `json:"outpoint"`
	AssetId   string   `json:"id"`
	ScriptKey string   `json:"script_key"`
	Amount    int
}

// FundVirtualPSBT sends a request to Tapd to fund a virtual PSBT using the invoice.
func FundVirtualPSBT(tapdHost, macaroon, invoice string, inputs PrevIds) (fundedPsbt *FundVirtualPSBTResponse, err error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/fund", tapdHost)

	// Disable TLS verification (for testing).
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Prepare the payload
	requestBody := map[string]interface{}{
		"raw": map[string]interface{}{
			"recipients": map[string]int{
				invoice: 0, // Use the invoice as the recipient
			},
			"inputs": inputs.Inputs,
		},
	}
	log.Println("Request body: ", requestBody)
	payloadBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
	}

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

	log.Printf("Response status: %s", resp.Status)

	// Parse the response
	err = json.NewDecoder(resp.Body).Decode(&fundedPsbt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return fundedPsbt, nil
}

type SignVirtualPSBTResponse struct {
	SignedPSBT string `json:"signed_psbt"`
}

// SignVirtualPSBT sends a request to Tapd to sign a virtual PSBT
func SignVirtualPSBT(tapdHost, macaroon, psbt string) (fundedPsbt *SignVirtualPSBTResponse, err error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/sign", tapdHost)

	// Disable TLS verification (for testing).
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Prepare the payload
	requestBody := map[string]interface{}{
		"funded_psbt": psbt,
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
		return nil, fmt.Errorf("tapd RPC error: %s", resp.Status)
	}

	// Parse the response into the AssetTransferResponse struct
	var assetTransfer AssetTransferResponse
	err = json.NewDecoder(resp.Body).Decode(&assetTransfer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &assetTransfer, nil
}
