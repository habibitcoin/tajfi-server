package tapd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendAssets sends a request to Tapd to send assets to the invoice.
// NOTE: THIS SHOULD ONLY BE USED WHEN DEMO MODE IS ENABLED.
func (c *tapdClient) SendAssets(tapdHost, macaroon, invoice string) (fundedPsbt *FundVirtualPSBTResponse, err error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/send", tapdHost)

	// Prepare the payload
	requestBody := map[string]interface{}{
		"tap_addrs": []string{invoice}, // Send invoice as a string array
	}
	payloadBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
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
		return nil, fmt.Errorf("tapd RPC error: %s", resp.Status)
	}

	// Parse the response
	err = json.NewDecoder(resp.Body).Decode(&fundedPsbt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return fundedPsbt, nil
}
