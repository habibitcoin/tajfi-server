package wallet

import (
	"net/http"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
)

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
	}
	if err := c.Validate(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed",
		})
	}

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
