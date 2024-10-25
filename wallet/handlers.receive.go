package wallet

import (
	"context"
	"log"
	"net/http"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
)

// Context key used in middleware to pass public key
type contextKey string

const pubKeyCtxKey = contextKey("public_key")

// Helper to extract public key from context
func getPublicKeyFromContext(ctx context.Context) (string, bool) {
	pubKey, ok := ctx.Value(pubKeyCtxKey).(string)
	return pubKey, ok
}

type RequestPayload struct {
	AssetID string `json:"asset_id" validate:"required"`
	Amount  int    `json:"amt" validate:"required"`
}

// StartSendAsset initiates a send transaction by calling the service.
func ReceiveAsset(c echo.Context) error {
	ctx := c.Request().Context()
	pubKey := ctx.Value("public_key").(string)
	// Extract public key from JWT

	// Bind and validate the request payload
	var payload RequestPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}
	/*if err := c.Validate(&payload); err != nil {
		log.Println(payload)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed",
		})
	}*/

	// Extract config from context
	cfg := config.GetConfig(c.Request().Context())

	// Call the service to start the send process
	params := ReceiveParams{
		PubKey:       pubKey,
		AssetID:      payload.AssetID,
		Amount:       payload.Amount,
		LNDHost:      cfg.LNDHost,
		LNMacaroon:   cfg.LNDMacaroon,
		TapdHost:     cfg.TapdHost,
		TapdMacaroon: cfg.TapdMacaroon,
	}

	response, err := Receive(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// this should move to the service.. we gotta clean up

	// if demo mode is enabled, and the amount is the same as demoAmount, pay it
	if cfg.DemoMode && payload.Amount == cfg.DemoAmount {
		_, err = tapd.SendAssets(cfg.DemoTapdHost, cfg.DemoTapdMacaroon, response["encoded"].(string))
		if err != nil {
			log.Println("Failed to send assets", err)
		}
	}

	return c.JSON(http.StatusOK, response)
}

// Handler for confirming a send transaction
func ConfirmSendAsset(c echo.Context) error {
	pubKey, ok := getPublicKeyFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Public key not found",
		})
	}

	// Use the public key to confirm the transaction
	// TODO: Add logic for confirming the transaction
	return c.JSON(http.StatusOK, map[string]string{
		"message":    "Send asset confirmed",
		"public_key": pubKey,
	})
}
