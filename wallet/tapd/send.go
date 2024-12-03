package tapd

import "fmt"

type PrevId struct {
	Outpoint    Outpoint `json:"outpoint"`
	AssetId     string   `json:"id"`
	ScriptKey   string   `json:"script_key"`
	Amount      int64    `json:"-"`
	InternalKey string   `json:"-"`
}

type Outpoint struct {
	Txid        string `json:"txid"`
	OutputIndex int    `json:"output_index"`
}

type PrevIds struct {
	Inputs []PrevId `json:"inputs,omitempty"`
}

type SignVirtualPSBTResponse struct {
	SignedPSBT string `json:"signed_psbt"`
}

type AnchorVirtualPSBTParams struct {
	VirtualPSBTs []string // Raw bytes for virtual PSBTs
	TapdHost     string
	Macaroon     string
}

type DecodeAddrResponse struct {
	Encoded   string `json:"encoded"`
	AssetID   string `json:"asset_id"`
	Amount    string `json:"amount"`
	ScriptKey string `json:"script_key"`
}

type FundVirtualPSBTResponse struct {
	FundedPSBT        string   `json:"funded_psbt"`
	ChangeOutputIndex int      `json:"change_output_index"`
	PassiveAssetPsbts []string `json:"passive_asset_psbts"`
	SighashHexToSign  string   `json:"sighash_hex_to_sign"`
}

func (c *tapdClient) DecodeAddr(tapdHost, macaroon, address string) (*DecodeAddrResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/addrs/decode", tapdHost)

	payload := map[string]string{"addr": address}
	var decodedAddr DecodeAddrResponse
	err := c.makePostRequest(url, macaroon, payload, &decodedAddr)
	return &decodedAddr, err
}

func (c *tapdClient) FundVirtualPSBT(tapdHost, macaroon, psbtOrInvoice string, inputs PrevIds) (*FundVirtualPSBTResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/fund", tapdHost)

	var payload map[string]interface{}

	if len(inputs.Inputs) > 0 {
		// Use raw format with recipients and inputs if no PSBT is provided.
		payload = map[string]interface{}{
			"raw": map[string]interface{}{
				"recipients": map[string]int{psbtOrInvoice: 0},
				"inputs":     inputs.Inputs,
			},
		}
	} else {
		// Use the PSBT format if a PSBT string is provided.
		payload = map[string]interface{}{
			"psbt": psbtOrInvoice,
		}
	}

	var fundedPsbt FundVirtualPSBTResponse
	err := c.makePostRequest(url, macaroon, payload, &fundedPsbt)
	return &fundedPsbt, err
}

func (c *tapdClient) SignVirtualPSBT(tapdHost, macaroon, psbt string) (*SignVirtualPSBTResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/sign", tapdHost)

	payload := map[string]interface{}{
		"funded_psbt": psbt,
	}

	var signedPsbt SignVirtualPSBTResponse
	err := c.makePostRequest(url, macaroon, payload, &signedPsbt)
	return &signedPsbt, err
}

// AnchorVirtualPSBT interacts with Tapd to anchor virtual PSBTs.
func (c *tapdClient) AnchorVirtualPSBT(params AnchorVirtualPSBTParams) (*AssetTransferResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/anchor", params.TapdHost)

	payload := map[string]interface{}{
		"virtual_psbts": params.VirtualPSBTs,
	}

	var assetTransfer AssetTransferResponse
	err := c.makePostRequest(url, params.Macaroon, payload, &assetTransfer)
	return &assetTransfer, err
}
