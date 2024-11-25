package tapd

import (
	"fmt"
)

// CreateInteractiveSendTemplateRequest represents the request for creating an interactive send template.
type CreateInteractiveSendTemplateRequest struct {
	AssetID           string `json:"asset_id,omitempty"`  // Hex-encoded asset ID.
	Amount            uint64 `json:"amount"`              // Amount of the asset to send.
	LockTime          uint64 `json:"lock_time,omitempty"` // Optional lock time.
	RelativeLockTime  uint64 `json:"relative_lock_time,omitempty"`
	OutputIndex       uint32 `json:"output_index,omitempty"`        // Anchor output index.
	ScriptKey         string `json:"script_key,omitempty"`          // Script key for receiving address.
	AnchorInternalKey string `json:"anchor_internal_key,omitempty"` // Internal key for the anchor transaction.
}

// CreateInteractiveSendTemplateResponse represents the response from creating an interactive send template.
type CreateInteractiveSendTemplateResponse struct {
	VirtualPSBT string `json:"psbt"` // Serialized virtual PSBT template.
}

// PrepareAnchoringTemplateRequest represents the request for preparing an anchoring template.
type PrepareAnchoringTemplateRequest struct {
	VirtualPSBT string `json:"virtual_psbt"` // Signed virtual PSBT.
	OutputAmt   int64  `json:"output_amt"`   // Amount in satoshis to send to the output address.
}

// PrepareAnchoringTemplateResponse represents the response for preparing an anchoring template.
type PrepareAnchoringTemplateResponse struct {
	AnchorPSBT string `json:"anchor_psbt"` // Serialized Bitcoin PSBT for anchoring.
}

// UpdateVirtualPsbtRequest represents the request for updating a virtual PSBT.
type UpdateVirtualPsbtRequest struct {
	VirtualPSBT string `json:"virtual_psbt"` // Serialized virtual PSBT to update.
	ScriptKey   string `json:"script_key"`   // New script key for the output.
	AnchorPSBT  string `json:"anchor_psbt"`  // Optional linked Bitcoin PSBT.
}

// UpdateVirtualPsbtResponse represents the response for updating a virtual PSBT.
type UpdateVirtualPsbtResponse struct {
	UpdatedVirtualPSBT string `json:"updated_virtual_psbt"` // Updated serialized virtual PSBT.
	UpdatedAnchorPSBT  string `json:"updated_anchor_psbt"`  // Updated Bitcoin PSBT (if provided).
}

// CreateInteractiveSendTemplate calls the RPC to create a virtual PSBT template for an interactive asset send.
func (c *tapdClient) CreateInteractiveSendTemplate(
	tapdHost, macaroon string,
	req CreateInteractiveSendTemplateRequest,
) (*CreateInteractiveSendTemplateResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/create-interactive-send-template", tapdHost)

	var response CreateInteractiveSendTemplateResponse
	err := c.makePostRequest(url, macaroon, req, &response)
	return &response, err
}

// PrepareAnchoringTemplate calls the RPC to prepare a Bitcoin PSBT for anchoring the given signed virtual PSBT.
func (c *tapdClient) PrepareAnchoringTemplate(
	tapdHost, macaroon string,
	req PrepareAnchoringTemplateRequest,
) (*PrepareAnchoringTemplateResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/prepare-anchoring-template", tapdHost)

	var response PrepareAnchoringTemplateResponse
	err := c.makePostRequest(url, macaroon, req, &response)
	return &response, err
}

// UpdateVirtualPSBT calls the RPC to update the script key of the output in the given virtual PSBT.
func (c *tapdClient) UpdateVirtualPSBT(
	tapdHost, macaroon string,
	req UpdateVirtualPsbtRequest,
) (*UpdateVirtualPsbtResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/update", tapdHost)

	var response UpdateVirtualPsbtResponse
	err := c.makePostRequest(url, macaroon, req, &response)
	return &response, err
}

// CommitVirtualPsbtsRequest represents the input for committing virtual PSBTs.
type CommitVirtualPsbtsRequest struct {
	VirtualPSBTs      []string `json:"virtual_psbts"`
	PassiveAssetPSBTs []string `json:"passive_asset_psbts"`
	AnchorPSBT        string   `json:"anchor_psbt"`
	AddChangeOutput   bool     `json:"add,omitempty"`
	TargetConf        uint32   `json:"target_conf,omitempty"`
	SatPerVByte       uint64   `json:"sat_per_vbyte,omitempty"`
}

// CommitVirtualPsbtsResponse represents the output of committing virtual PSBTs.
type CommitVirtualPsbtsResponse struct {
	AnchorPSBT        string     `json:"anchor_psbt"`
	VirtualPSBTs      []string   `json:"virtual_psbts"`
	PassiveAssetPSBTs []string   `json:"passive_asset_psbts"`
	ChangeOutputIndex int32      `json:"change_output_index"`
	LndLockedUtxos    []Outpoint `json:"lnd_locked_utxos"`
}

// CommitVirtualPsbts commits virtual PSBTs to the anchor PSBT.
func (c *tapdClient) CommitVirtualPsbts(tapdHost, macaroon string, req CommitVirtualPsbtsRequest) (*CommitVirtualPsbtsResponse, error) {
	url := fmt.Sprintf("https://%s/v1/taproot-assets/wallet/virtual-psbt/commit", tapdHost)

	var response CommitVirtualPsbtsResponse
	err := c.makePostRequest(url, macaroon, req, &response)
	return &response, err
}
