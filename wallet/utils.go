package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"tajfi-server/wallet/tapd"

	btcec "github.com/btcsuite/btcd/btcec/v2"
)

// CompressPubKey converts a 32-byte X-coordinate into a 33-byte compressed public key.
func CompressPubKey(pubKeyHex string) (string, error) {
	// Decode the provided X-coordinate from hex
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", errors.New("invalid public key hex")
	}
	if len(pubKeyBytes) != 32 {
		return "", errors.New("invalid length: expected 32 bytes")
	}

	// Use secp256k1 curve
	curve := btcec.S256()

	// Convert the X-coordinate to a big.Int
	x := new(big.Int).SetBytes(pubKeyBytes)

	// Compute the Y-coordinate using the curve equation: y² = x³ + 7 (mod p)
	ySquared := new(big.Int).Exp(x, big.NewInt(3), curve.P) // x³
	ySquared.Add(ySquared, curve.B)                         // x³ + 7
	y := new(big.Int).ModSqrt(ySquared, curve.P)            // √(x³ + 7)

	if y == nil {
		return "", errors.New("no valid Y-coordinate found")
	}

	// Check if the Y-coordinate is even or odd
	prefix := byte(0x02) // Even prefix
	if y.Bit(0) == 1 {
		prefix = 0x03 // Odd prefix
	}

	// Create the compressed public key: [prefix || X-coordinate]
	compressedKey := append([]byte{prefix}, pubKeyBytes...)

	return hex.EncodeToString(compressedKey), nil
}

// Helper function to read the override signature or sighash from a file.
func ReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, return an empty string.
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// Save sigHash to a file named sighash.hex, clearing any previous contents.
func WriteSignatureToFile(filename string, signatureHex string) error {
	// Open the file with write-only mode, create it if it doesn't exist, and truncate it to clear existing content.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open/create file: %w", err)
	}
	defer file.Close()

	// Write the sigHash as a hex string to the file.
	_, err = file.WriteString(signatureHex)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("Signature successfully written to " + filename)
	return nil
}

func FilterOwnedUtxos(utxos *tapd.GetUtxosResponse, pubKey string) (ownedUtxos tapd.PrevIds) {
	for _, utxo := range utxos.ManagedUtxos {
		for _, asset := range utxo.Assets {
			if asset.ScriptKey == ("02" + pubKey) {
				txid, vout, err := parseOutPoint(utxo.Outpoint)
				if err != nil {
					fmt.Println("Error parsing outpoint:", err)
					continue
				}
				ownedUtxos.Inputs = append(ownedUtxos.Inputs, tapd.PrevId{
					Outpoint:  tapd.Outpoint{Txid: txid, OutputIndex: vout},
					AssetId:   asset.AssetGenesis.AssetID,
					ScriptKey: asset.ScriptKey,
				})
				break
			}
		}
	}
	return ownedUtxos
}

func parseOutPoint(outPoint string) (txid string, vout int, err error) {
	parts := strings.Split(outPoint, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid out_point format")
	}

	txid = parts[0]
	_, err = fmt.Sscanf(parts[1], "%d", &vout)
	if err != nil {
		return "", 0, fmt.Errorf("invalid vout format: %v", err)
	}

	return txid, vout, nil
}
