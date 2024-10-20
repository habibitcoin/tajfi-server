package lnd

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type InternalKeyResponse struct {
	RawKeyBytes string `json:"raw_key_bytes"`
	KeyLoc      struct {
		KeyFamily int `json:"key_family"`
		KeyIndex  int `json:"key_index"`
	} `json:"key_loc"`
}

// GetInternalKey retrieves an internal key from the LND node.
func GetInternalKey(lndHost, macaroon string) (*InternalKeyResponse, error) {
	url := fmt.Sprintf("https://%s/v2/wallet/key/next", lndHost)

	// Create request payload
	payload := map[string]interface{}{
		"key_family": 212,
	}
	payloadBytes, _ := json.Marshal(payload)

	// Disable TLS verification for simplicity (use with caution!)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get internal key: %s", body)
	}

	// Parse response
	var keyResponse InternalKeyResponse
	err = json.NewDecoder(resp.Body).Decode(&keyResponse)
	if err != nil {
		return nil, err
	}

	// Decode raw key bytes from base64
	rawBytes, err := base64.StdEncoding.DecodeString(keyResponse.RawKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw key bytes: %v", err)
	}
	keyResponse.RawKeyBytes = fmt.Sprintf("%x", rawBytes) // Hexify

	return &keyResponse, nil
}
