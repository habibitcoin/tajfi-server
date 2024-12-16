package wallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"tajfi-server/wallet/tapd"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
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

func FilterOwnedUtxos(utxos *tapd.GetUtxosResponse, pubKey string, assetId string) (ownedUtxos tapd.PrevIds) {
	for _, utxo := range utxos.ManagedUtxos {
		for _, asset := range utxo.Assets {
			if asset.ScriptKey == ("02"+pubKey) && asset.AssetGenesis.AssetID == assetId {
				txid, vout, err := parseOutPoint(utxo.Outpoint)
				if err != nil {
					fmt.Println("Error parsing outpoint:", err)
					continue
				}
				amount, err := strconv.ParseInt(asset.Amount, 10, 64)
				if err != nil {
					fmt.Println("Error converting amount:", err)
					continue
				}
				ownedUtxos.Inputs = append(ownedUtxos.Inputs, tapd.PrevId{
					Outpoint:    tapd.Outpoint{Txid: txid, OutputIndex: vout},
					AssetId:     asset.AssetGenesis.AssetID,
					ScriptKey:   asset.ScriptKey,
					Amount:      amount,
					InternalKey: utxo.InternalKey,
				})
				break
			}
		}
	}
	// Sort the ownedUtxos.Inputs slice by Amount in descending order
	sort.Slice(ownedUtxos.Inputs, func(i, j int) bool {
		return ownedUtxos.Inputs[i].Amount > ownedUtxos.Inputs[j].Amount
	})
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

// WriteOrderToFile writes an order to a JSON file.
func WriteOrderToFile(order Order, filename string) error {
	data, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal order to JSON: %w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write order to file: %w", err)
	}

	return nil
}

// ReadOrderFromFile reads an order from a JSON file.
func ReadOrderFromFile(filename string) (Order, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Order{}, fmt.Errorf("failed to read order file: %w", err)
	}

	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return Order{}, fmt.Errorf("failed to unmarshal order JSON: %w", err)
	}

	return order, nil
}

// FetchOutpointFromPsbt extracts the Txid and OutputIndex of the first input from a PSBT.
func FetchOutpointFromPsbt(psbtHex string) (tapd.Outpoint, error) {
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		return tapd.Outpoint{}, fmt.Errorf("failed to decode PSBT hex: %w", err)
	}

	parsedPsbt, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return tapd.Outpoint{}, fmt.Errorf("failed to parse PSBT: %w", err)
	}

	if len(parsedPsbt.UnsignedTx.TxIn) == 0 {
		return tapd.Outpoint{}, fmt.Errorf("no inputs in PSBT")
	}

	input := parsedPsbt.UnsignedTx.TxIn[0]
	txid := hex.EncodeToString(reverse(input.PreviousOutPoint.Hash[:])) // Reverse for correct Txid order
	return tapd.Outpoint{
		Txid:        txid,
		OutputIndex: int(input.PreviousOutPoint.Index),
	}, nil
}

// reverse reverses a slice of bytes.
func reverse(input []byte) []byte {
	reversed := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		reversed[i] = input[len(input)-1-i]
	}
	return reversed
}

// generatePkScript generates the PkScript (locking script) for a given Bitcoin address.
func GeneratePkScript(address string) ([]byte, error) {
	// Decode the address into a btcutil.Address object.
	addr, err := btcutil.DecodeAddress(address, &chaincfg.SigNetParams) // Adjust network params as needed.
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}

	// Generate the PkScript for the address.
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create PkScript: %w", err)
	}

	return pkScript, nil
}
