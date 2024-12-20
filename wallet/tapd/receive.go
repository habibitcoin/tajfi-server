package tapd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"tajfi-server/wallet/lnd"
)

type NewAddressPayload struct {
	AssetID     string                  `json:"asset_id"`
	Amt         int                     `json:"amt"`
	ScriptKey   map[string]interface{}  `json:"script_key"`
	InternalKey lnd.InternalKeyResponse `json:"internal_key"`
}

// CallNewAddress sends the payload to the Tapd NewAddress RPC.
func (c *tapdClient) CallNewAddress(tapdHost, macaroon string, payload NewAddressPayload) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/addrs", tapdHost)

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tapd RPC error: %s", resp.Status)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
