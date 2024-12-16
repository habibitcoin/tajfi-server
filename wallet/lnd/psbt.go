package lnd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

// SignPsbtRequest represents the payload for signing a PSBT.
type SignPsbtRequest struct {
	FundedPsbt string `json:"funded_psbt"`
}

// SignPsbtResponse represents the response after signing a PSBT.
type SignPsbtResponse struct {
	SignedPsbt   string   `json:"signed_psbt"`
	SignedInputs []uint32 `json:"signed_inputs,omitempty"`
}

// SignPsbt sends a request to the LND node to sign a PSBT.
func SignPsbt(lndHost, macaroon string, req SignPsbtRequest) (*SignPsbtResponse, error) {
	url := fmt.Sprintf("https://%s/v2/wallet/psbt/sign", lndHost)

	payloadBytes, _ := json.Marshal(req)
	reqBody := bytes.NewBuffer(payloadBytes)

	// Disable TLS verification for simplicity (use with caution!)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Grpc-Metadata-macaroon", macaroon)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LND SignPsbt RPC error: %s", resp.Status)
	}

	var signResp SignPsbtResponse
	if err := json.NewDecoder(resp.Body).Decode(&signResp); err != nil {
		return nil, err
	}

	return &signResp, nil
}

// FinalizePsbt sends a request to the LND node to sign a PSBT.
func FinalizePsbt(lndHost, macaroon string, req SignPsbtRequest) (*SignPsbtResponse, error) {
	url := fmt.Sprintf("https://%s/v2/wallet/psbt/finalize", lndHost)

	payloadBytes, _ := json.Marshal(req)
	reqBody := bytes.NewBuffer(payloadBytes)

	// Disable TLS verification for simplicity (use with caution!)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Grpc-Metadata-macaroon", macaroon)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LND SignPsbt RPC error: %s", resp.Status)
	}

	var signResp SignPsbtResponse
	if err := json.NewDecoder(resp.Body).Decode(&signResp); err != nil {
		return nil, err
	}

	return &signResp, nil
}
