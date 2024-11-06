package wallet

import (
	"strconv"
	"strings"
	"tajfi-server/wallet/tapd"
)

func GetTransfersResponse(tapdTransfers tapd.AssetTransfersResponse, pubKey string) (transfers []Transfer) {
	pubKey = "02" + pubKey // normalize the public key
	for _, tapdTransfer := range tapdTransfers.Transfers {
		var transfer Transfer
		transfer.Timestamp = tapdTransfer.TransferTimestamp
		transfer.Height = tapdTransfer.AnchorTxHeightHint
		transfer.Txid = strings.Split(tapdTransfer.Outputs[0].Anchor.Outpoint, ":")[0]

		if len(tapdTransfer.Inputs) > 0 {
			transfer.AssetID = tapdTransfer.Inputs[0].AssetID
		}

		var sentAmount, receivedAmount uint64

		for _, input := range tapdTransfer.Inputs {
			if input.ScriptKey == pubKey {
				amount, err := strconv.ParseUint(input.Amount, 10, 64)
				if err == nil {
					sentAmount += amount
				}
			}
		}

		for _, output := range tapdTransfer.Outputs {
			if output.ScriptKey == pubKey {
				amount, err := strconv.ParseUint(output.Amount, 10, 64)
				if err == nil {
					receivedAmount += amount
				}
			}
		}

		if receivedAmount == 0 && sentAmount == 0 {
			continue
		}

		if receivedAmount > sentAmount {
			transfer.Type = "receive"
			transfer.Amount = receivedAmount - sentAmount
		} else {
			transfer.Type = "send"
			transfer.Amount = sentAmount - receivedAmount
		}

		transfers = append(transfers, transfer)
	}

	return transfers
}

func GenerateNewAddress() string {
	// Logic to generate a Taproot address
	return "bc1qxyz..."
}

func GetBalance() int {
	// Logic to fetch wallet balance
	return 100
}
