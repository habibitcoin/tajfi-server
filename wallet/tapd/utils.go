package tapd

import "fmt"

func (c *tapdClient) SendAssets(tapdHost, macaroon, invoice string) (*FundVirtualPSBTResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/send", tapdHost)

	payload := map[string]interface{}{
		"tap_addrs": []string{invoice},
	}

	var fundedPsbt FundVirtualPSBTResponse
	err := c.makePostRequest(url, macaroon, payload, &fundedPsbt)
	return &fundedPsbt, err
}
