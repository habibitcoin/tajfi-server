package wallet

import (
	"net/http"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
)

type SendDecodePayload struct {
	Address string `json:"address" validate:"required"`
}

// DecodeAddress handles decoding of Taproot Asset addresses
func DecodeAddress(c echo.Context) error {
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
	decoded, err := tapd.DecodeAddr(cfg.TapdHost, cfg.TapdMacaroon, payload.Address)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode address")
	}

	return c.JSON(http.StatusOK, decoded)
}

// SendStartPayload defines the request payload structure for /send/start.
type SendStartPayload struct {
	Invoice string `json:"invoice" validate:"required"`
}

// SendStart initiates the send transaction by calling Tapd and returning a vPSBT
func SendStart(c echo.Context) error {
	// Parse the request payload
	var payload SendStartPayload
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

	// Call the Tapd service to fund the PSBT
	fundedPsbt, err := tapd.FundVirtualPSBT(cfg.TapdHost, cfg.TapdMacaroon, payload.Invoice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, fundedPsbt)
}

// SendCompletePayload defines the request payload structure for /send/complete.
type SendCompletePayload struct {
	PSBT string `json:"psbt" validate:"required"`
}

// SendComplete completes the send transaction by calling Tapd and returning the transaction ID
func SendComplete(c echo.Context) error {
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

	params := tapd.AnchorVirtualPSBTParams{
		VirtualPSBTs: []string{payload.PSBT},
		TapdHost:     cfg.TapdHost,
		Macaroon:     cfg.TapdMacaroon,
	}

	// Call the Tapd service to fund the PSBT
	fundedPsbt, err := tapd.AnchorVirtualPSBT(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, fundedPsbt)
}
