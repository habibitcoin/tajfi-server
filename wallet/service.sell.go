package wallet

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"tajfi-server/wallet/lnd"
	"tajfi-server/wallet/tapd"

	"github.com/btcsuite/btcd/btcutil/psbt"
)

// SellStartParams contains the parameters for starting a sell process.
type SellStartParams struct {
	PubKey       string
	AssetID      string
	AmountToSell int64
	TapdHost     string
	TapdMacaroon string
}

// SellCompleteParams contains the parameters for completing a sell process.
type SellCompleteParams struct {
	Signature           string
	PSBT                string
	SighashHex          string
	SignatureHex        string
	AmountSatsToReceive int64
	ReceiveBTCAddress   string
	TapdHost            string
	TapdMacaroon        string
}

// StartSellService handles the creation of a virtual PSBT for selling assets.
func StartSellService(params SellStartParams, tapdClient tapd.TapdClientInterface) (*SellStartResponse, error) {
	cfg := tapdClient.GetClientConfig()

	// First, we need to check if the user has a UTXO that can be used to fund the sell.
	// NOTE: Only 1 UTXO is supported currently.
	utxos, err := tapdClient.GetUtxos(cfg.TapdHost, cfg.TapdMacaroon)
	if err != nil {
		return nil, fmt.Errorf("failed to get utxos: %w", err)
	}
	eligibleUtxos := FilterOwnedUtxos(utxos, params.PubKey, params.AssetID)

	if len(eligibleUtxos.Inputs) == 0 {
		return nil, fmt.Errorf("no eligible UTXOs found to sell for assetID %s", params.AssetID)
	}
	var (
		selectedUtxo tapd.PrevId
		validAmounts []int64
	)
	for _, utxo := range eligibleUtxos.Inputs {
		if utxo.Amount == params.AmountToSell {
			selectedUtxo = utxo
			break
		}
		validAmounts = append(validAmounts, utxo.Amount)
	}
	if selectedUtxo == (tapd.PrevId{}) {
		return nil, fmt.Errorf("no UTXO found with exact amount %d for assetID %s. eligible amounts are: %v", params.AmountToSell, params.AssetID, validAmounts)
	}

	// Step 1: Call CreateInteractiveSendTemplate
	templateReq := tapd.CreateInteractiveSendTemplateRequest{
		AssetID:           params.AssetID,
		Amount:            uint64(params.AmountToSell),
		ScriptKey:         "02" + params.PubKey,
		AnchorInternalKey: eligibleUtxos.Inputs[0].InternalKey,
	}
	templateResp, err := tapdClient.CreateInteractiveSendTemplate(params.TapdHost, params.TapdMacaroon, templateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create interactive send template: %w", err)
	}

	// Step 2: Call FundVirtualPSBT
	fundResp, err := tapdClient.FundVirtualPSBT(params.TapdHost, params.TapdMacaroon, templateResp.VirtualPSBT, tapd.PrevIds{Inputs: []tapd.PrevId{selectedUtxo}}, false)
	if err != nil {
		return nil, fmt.Errorf("failed to fund virtual PSBT: %w", err)
	}

	// Call the modified Tapd service to write our sighash.hex file
	_, err = tapdClient.SignVirtualPSBT(params.TapdHost, params.TapdMacaroon, fundResp.FundedPSBT)
	// if err != nil { we dont really care if this fails, it is expected }

	// Fetch the hex file now
	taprootSigsDir := tapdClient.GetClientConfig().TaprootSigsDir
	sighash, err := ReadFile(taprootSigsDir + "sighash.hex")
	if err != nil {
		return nil, fmt.Errorf("failed to read sighash.hex: %w", err)
	}
	fundResp.SighashHexToSign = sighash

	// Use the selected UTXO's Outpoint for the filename and order.
	filename := fmt.Sprintf("./orders/pending/%s_%d.json",
		selectedUtxo.Outpoint.Txid,
		selectedUtxo.Outpoint.OutputIndex,
	)

	order := Order{
		AssetID:             params.AssetID,
		AmountToSell:        params.AmountToSell,
		AmountSatsToReceive: 0, // To be updated in CompleteSellService
		Outpoint: tapd.Outpoint{
			Txid:        selectedUtxo.Outpoint.Txid,
			OutputIndex: selectedUtxo.Outpoint.OutputIndex,
		},
	}

	// Save the order to a JSON file.
	if err := WriteOrderToFile(order, filename); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Step 3: Return the unsigned PSBT and sighash
	return &SellStartResponse{
		FundVirtualPSBTResponse: *fundResp,
	}, nil
}

// CompleteSellService handles the completion of the sell process by signing the virtual PSBT and preparing the anchor.
func CompleteSellService(params SellCompleteParams, tapdClient tapd.TapdClientInterface) (*SellCompleteResponse, error) {
	// Write the override signature to the file that tapd recognizes.
	cfg := tapdClient.GetClientConfig()
	if err := WriteSignatureToFile(cfg.TaprootSigsDir+"signature.hex", params.SignatureHex); err != nil {
		return nil, fmt.Errorf("failed to write signature to file: %w", err)
	}

	// Step 1: SignVirtualPSBT.
	signResp, err := tapdClient.SignVirtualPSBT(params.TapdHost, params.TapdMacaroon, params.PSBT)
	if err != nil {
		return nil, fmt.Errorf("failed to sign virtual PSBT: %w", err)
	}
	log.Println("Signed PSBT successfully")

	// now delete the signature.hex and sighash.hex files
	err = os.Remove(cfg.TaprootSigsDir + "signature.hex")
	if err != nil {
		log.Println("Failed to delete signature.hex file", err)
	}
	err = os.Remove(cfg.TaprootSigsDir + "sighash.hex")
	if err != nil {
		log.Println("Failed to delete sighash.hex file", err)
	}

	// Step 2: Prepare the anchoring template.
	anchorReq := tapd.PrepareAnchoringTemplateRequest{
		VirtualPSBT: signResp.SignedPSBT,
		OutputAmt:   params.AmountSatsToReceive,
	}
	anchorResp, err := tapdClient.PrepareAnchoringTemplate(params.TapdHost, params.TapdMacaroon, anchorReq)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare anchoring template: %w", err)
	}
	log.Println("Prepared anchoring template successfully")

	// Step 3: CommitVirtualPSBTs.
	commitReq := tapd.CommitVirtualPsbtsRequest{
		VirtualPSBTs:    []string{signResp.SignedPSBT},
		AnchorPSBT:      anchorResp.AnchorPSBT,
		AddChangeOutput: true,
		TargetConf:      6,
	}
	commitResp, err := tapdClient.CommitVirtualPsbts(params.TapdHost, params.TapdMacaroon, commitReq)
	if err != nil {
		return nil, fmt.Errorf("failed to commit virtual PSBTs: %w", err)
	}

	log.Println("Committed virtual PSBTs successfully")
	log.Printf("Anchor PSBT: %s\n", commitResp.AnchorPSBT)
	log.Printf("Virtual PSBT: %s\n", commitResp.VirtualPSBTs)

	// Step 4: Modify PSBT to strip funding inputs/outputs.
	modifiedPSBT := commitResp.AnchorPSBT
	modifiedPSBT, err = stripFundingInputsAndOutputs(modifiedPSBT, params.ReceiveBTCAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to strip funding inputs/outputs: %w", err)
	}

	log.Println("Modified PSBT successfully")
	log.Printf("Modified PSBT: %s\n", modifiedPSBT)

	// Step 5: Sign the modified PSBT with LND.
	// Convert the hex string to a byte slice
	modifiedPSBTBytes, err := hex.DecodeString(modifiedPSBT)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modified PSBT from hex: %w", err)
	}
	signReq := lnd.SignPsbtRequest{FundedPsbt: base64.StdEncoding.EncodeToString(modifiedPSBTBytes)}
	lndSignResp, err := lnd.SignPsbt(cfg.LNDHost, cfg.LNDMacaroon, signReq)
	if err != nil {
		return nil, fmt.Errorf("failed to sign PSBT with LND: %w", err)
	}

	// Return the finalized PSBTs.
	signedPsbtBytes, err := base64.StdEncoding.DecodeString(lndSignResp.SignedPsbt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 signed PSBT: %w", err)
	}
	signedPsbtHex := hex.EncodeToString(signedPsbtBytes)

	// Locate the order file using LndLockedUtxos.

	inputUtxo, err := FetchOutpointFromPsbt(modifiedPSBT)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch outpoint from PSBT: %w", err)
	}
	filename := fmt.Sprintf("./orders/pending/%s_%d.json",
		inputUtxo.Txid,
		inputUtxo.OutputIndex,
	)
	for _, lockedUtxo := range commitResp.LndLockedUtxos {
		log.Println("Locked UTXO:", lockedUtxo)
		// TODO: Should immediately release the lock for these UTXOs.
	}

	// Read the existing order.
	order, err := ReadOrderFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read order file: %w", err)
	}

	// Update the order with additional PSBT details.
	order.VirtualPSBT = signResp.SignedPSBT
	order.AnchorPSBT = anchorResp.AnchorPSBT
	order.PassiveAssetPSBTs = commitResp.VirtualPSBTs
	order.AmountSatsToReceive = params.AmountSatsToReceive
	// remove /pending from the filename
	finalFileName := strings.Replace(filename, "/pending", "", 1)
	// Save the updated order.
	if err := WriteOrderToFile(order, finalFileName); err != nil {
		return nil, fmt.Errorf("failed to update order file: %w", err)
	}
	// delete pending order
	err = os.Remove(filename)
	if err != nil {
		log.Println("Failed to delete pending order file", err)
	}

	return &SellCompleteResponse{
		SignedVirtualPSBT:  signResp.SignedPSBT,
		ModifiedAnchorPSBT: signedPsbtHex,
	}, nil
}

// stripFundingInputsAndOutputs removes unnecessary funding inputs and outputs from a PSBT.
// we only support 1-in-2-out PSBTs for now.
func stripFundingInputsAndOutputs(psbtHex string, destinationAddress string) (string, error) {
	// Decode the hex string into raw PSBT bytes.
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode PSBT hex: %w", err)
	}

	// Deserialize the PSBT using a bytes.Reader.
	btcpsbt, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize PSBT: %w", err)
	}

	// Modify inputs: Remove the funding input at index 1 (keep input 0 and input 2).
	if len(btcpsbt.Inputs) > 1 {
		btcpsbt.Inputs = append(btcpsbt.Inputs[:1], btcpsbt.Inputs[2:]...)
		btcpsbt.UnsignedTx.TxIn = append(btcpsbt.UnsignedTx.TxIn[:1], btcpsbt.UnsignedTx.TxIn[2:]...)
	}

	// Modify outputs: Keep only the first two outputs.
	if len(btcpsbt.Outputs) > 2 {
		btcpsbt.Outputs = btcpsbt.Outputs[:2]
		btcpsbt.UnsignedTx.TxOut = btcpsbt.UnsignedTx.TxOut[:2]
	}
	if len(btcpsbt.UnsignedTx.TxOut) >= 1 {
		// Generate the PkScript for the input address.
		if destinationAddress != "" {
			pkScript, err := GeneratePkScript(destinationAddress)
			if err != nil {
				return "", fmt.Errorf("failed to generate PkScript for address: %w", err)
			}

			// Update the first output's PkScript.
			btcpsbt.UnsignedTx.TxOut[0].PkScript = pkScript
		}
		btcpsbt.UnsignedTx.TxOut[1].Value = 1000 // Set the amount to 1000 sats.
	}

	// Serialize the modified PSBT back into raw bytes.
	var buf bytes.Buffer
	err = btcpsbt.Serialize(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to serialize modified PSBT: %w", err)
	}

	// Encode the modified PSBT back to hex format.
	return hex.EncodeToString(buf.Bytes()), nil
}
