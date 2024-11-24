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
		transfer.Status = "confirmed"
		if tapdTransfer.AnchorTxBlockHash.Hash == "" {
			transfer.Status = "unconfirmed"
		}
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
			transfer.ChangeAmount = receivedAmount
		}

		transfers = append(transfers, transfer)
	}

	return transfers
}

func UpdateBalancesWithUnconfirmed(
	balances *tapd.WalletBalancesResponse,
	transfers []Transfer,
) {
	for _, transfer := range transfers {
		// Skip confirmed transfers
		if transfer.Status == "confirmed" {
			continue
		}

		var (
			genesisID = transfer.AssetID
			amount    = transfer.Amount
		)

		if transfer.Type == "send" {
			// For sends, we need to add the change amount back
			amount = transfer.ChangeAmount
		}

		// Update unconfirmed balance for the relevant asset
		if balance, exists := balances.AssetBalances[genesisID]; exists {
			// Update the unconfirmed balance
			existingUnconfirmed, _ := strconv.Atoi(balance.UnconfirmedBalance)
			balance.UnconfirmedBalance = strconv.Itoa(existingUnconfirmed + int(amount))

			// Update the confirmed balance
			existingConfirmed, _ := strconv.Atoi(balance.Balance)
			balance.Balance = strconv.Itoa(existingConfirmed + int(amount))

			// Reassign the updated balance back to the map
			balances.AssetBalances[genesisID] = balance
		} else {
			amountStr := strconv.Itoa(int(amount))
			// Create a new entry for unconfirmed balance
			balances.AssetBalances[genesisID] = tapd.AssetBalance{
				AssetGenesis:       tapd.AssetGenesis{AssetID: genesisID},
				Balance:            amountStr,
				UnconfirmedBalance: amountStr,
			}
		}
	}
}

func GenerateNewAddress() string {
	// Logic to generate a Taproot address
	return "bc1qxyz..."
}

func GetBalance() int {
	// Logic to fetch wallet balance
	return 100
}
