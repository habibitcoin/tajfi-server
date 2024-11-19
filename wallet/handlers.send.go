package wallet

import (
	"log"
	"net/http"
	"os"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
)

type SendDecodePayload struct {
	Address string `json:"address" validate:"required"`
}

// DecodeAddress handles decoding of Taproot Asset addresses
func DecodeAddress(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload SendDecodePayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		}
		/*
			if err := c.Validate(&payload); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Validation failed",
				})
			}*/

		// Extract config from context
		cfg := config.GetConfig(c.Request().Context())

		// Call the tapd package's DecodeAddr method
		decoded, err := tapdClient.DecodeAddr(cfg.TapdHost, cfg.TapdMacaroon, payload.Address)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode address")
		}

		return c.JSON(http.StatusOK, decoded)
	}
}

// SendStartPayload defines the request payload structure for /send/start.
type SendStartPayload struct {
	Invoice string `json:"invoice" validate:"required"`
}

// SendStart initiates the send transaction by calling Tapd and returning a vPSBT
func SendStart(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request payload
		var payload SendStartPayload
		ctx := c.Request().Context()
		pubKey := ctx.Value("public_key").(string)
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload",
			})
		} /*
			if err := c.Validate(&payload); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Validation failed",
				})
			}*/

		// Extract config from context
		cfg := config.GetConfig(c.Request().Context())

		utxos, err := tapdClient.GetUtxos(cfg.TapdHost, cfg.TapdMacaroon)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch balances from tapd: "+err.Error())
		}

		// Decode the address
		decoded, err := tapdClient.DecodeAddr(cfg.TapdHost, cfg.TapdMacaroon, payload.Invoice)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode address")
		}

		myUtxos := FilterOwnedUtxos(utxos, pubKey, decoded.AssetID)
		log.Printf("Found %d UTXOs for pubkey %s", len(myUtxos.Inputs), pubKey)
		// Call the Tapd service to fund the PSBT
		fundedPsbt, err := tapdClient.FundVirtualPSBT(cfg.TapdHost, cfg.TapdMacaroon, payload.Invoice, tapd.PrevIds{myUtxos.Inputs})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		// Call the modified Tapd service to write our sighash.hex file
		_, err = tapdClient.SignVirtualPSBT(cfg.TapdHost, cfg.TapdMacaroon, fundedPsbt.FundedPSBT)
		// if err != nil { we dont really care if this fails, it is expected }

		// Fetch the hex file now
		taprootSigsDir := cfg.TaprootSigsDir
		sighash, err := ReadFile(taprootSigsDir + "sighash.hex")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		fundedPsbt.SighashHexToSign = sighash

		return c.JSON(http.StatusOK, fundedPsbt)
	}
}

// SendCompletePayload defines the request payload structure for /send/complete.
type SendCompletePayload struct {
	PSBT         string `json:"psbt" validate:"required"`
	SignatureHex string `json:"signature_hex" validate:"required"`
}

// SendComplete completes the send transaction by calling Tapd and returning the transaction ID
func SendComplete(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request payload
		var payload SendCompletePayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload",
			})
		} /*
			if err := c.Validate(&payload); err != nil	{
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Validation failed",
				})
			}*/
		// Extract config from context
		cfg := config.GetConfig(c.Request().Context())

		// Write the override signature to the file that tapd recognizes
		if err := WriteSignatureToFile(cfg.TaprootSigsDir+"signature.hex", payload.SignatureHex); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		// Call the Tapd service to sign the PSBT, knowing it will use the signature.hex file
		signedPsbt, err := tapdClient.SignVirtualPSBT(cfg.TapdHost, cfg.TapdMacaroon, payload.PSBT)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
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

		params := tapd.AnchorVirtualPSBTParams{
			VirtualPSBTs: []string{signedPsbt.SignedPSBT},
			TapdHost:     cfg.TapdHost,
			Macaroon:     cfg.TapdMacaroon,
		}

		// Call the Tapd service to fund the PSBT
		fundedPsbt, err := tapdClient.AnchorVirtualPSBT(params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, fundedPsbt)
	}
}
