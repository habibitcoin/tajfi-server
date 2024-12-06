package tapd

import (
	"fmt"
	"tajfi-server/wallet/lnd"
)

type NewAddressPayload struct {
	AssetID     string                  `json:"asset_id"`
	Amt         int                     `json:"amt"`
	ScriptKey   map[string]interface{}  `json:"script_key"`
	InternalKey lnd.InternalKeyResponse `json:"internal_key"`
}

func (c *tapdClient) CallNewAddress(tapdHost, macaroon string, payload NewAddressPayload) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/addrs", tapdHost)

	var result map[string]interface{}
	err := c.makePostRequest(url, macaroon, payload, &result)
	return result, err
}
