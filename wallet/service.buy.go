package wallet

import (
	"fmt"
	"log"
	"tajfi-server/wallet/lnd"
	"tajfi-server/wallet/tapd"
)

type BuyStartParams struct {
	PSBT         string
	AnchorPSBT   string
	PubKey       string
	TapdHost     string
	TapdMacaroon string
}

type BuyCompleteParams struct {
	PSBT            string
	AnchorPSBT      string
	SighashHex      string
	SignatureHex    string
	AmountSatsToPay int64
	TapdHost        string
	TapdMacaroon    string
}

// StartBuyService updates the virtual PSBT for the buy process.
func StartBuyService(params BuyStartParams, tapdClient tapd.TapdClientInterface) (*tapd.UpdateVirtualPsbtResponse, error) {
	log.Printf("Using anchor PSBT: %s", params.AnchorPSBT)
	updateReq := tapd.UpdateVirtualPsbtRequest{
		VirtualPSBT: params.PSBT,
		ScriptKey:   "02" + params.PubKey,
		AnchorPSBT:  params.AnchorPSBT,
	}
	updateResp, err := tapdClient.UpdateVirtualPSBT(params.TapdHost, params.TapdMacaroon, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update virtual PSBT: %w", err)
	}
	log.Printf("Updated virtual PSBT: %s", updateResp.UpdatedVirtualPSBT)
	log.Printf("Anchor PSBT: %s", updateResp.UpdatedAnchorPSBT)

	return updateResp, nil
}

// CompleteBuyService finalizes the buy process by signing and finalizing the PSBT.
func CompleteBuyService(params BuyCompleteParams, tapdClient tapd.TapdClientInterface) (*SellCompleteResponse, error) {
	cfg := tapdClient.GetClientConfig()

	// Commit the PSBTs.
	commitReq := tapd.CommitVirtualPsbtsRequest{
		VirtualPSBTs:    []string{params.PSBT},
		AnchorPSBT:      params.AnchorPSBT,
		AddChangeOutput: true,
		TargetConf:      6,
	}
	commitResp, err := tapdClient.CommitVirtualPsbts(params.TapdHost, params.TapdMacaroon, commitReq)
	if err != nil {
		return nil, fmt.Errorf("failed to commit virtual PSBTs: %w", err)
	}
	log.Printf("Committed PSBTs: %s", commitResp.VirtualPSBTs)
	log.Printf("Anchor PSBT: %s", commitResp.AnchorPSBT)

	// Fund and sign the committed PSBT using LND.
	signReq := lnd.SignPsbtRequest{FundedPsbt: string(commitResp.AnchorPSBT)}
	finalizedPSBT, err := lnd.SignPsbt(cfg.LNDHost, cfg.LNDMacaroon, signReq)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize PSBT: %w", err)
	}

	return &SellCompleteResponse{
		SignedVirtualPSBT:  commitResp.VirtualPSBTs[0],
		ModifiedAnchorPSBT: finalizedPSBT.SignedPsbt,
	}, nil
}
