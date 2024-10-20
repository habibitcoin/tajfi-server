package wallet

import (
	"encoding/hex"
	"errors"
	"math/big"

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
