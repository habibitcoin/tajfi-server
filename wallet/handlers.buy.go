package wallet

import (
	"net/http"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
)

type BuyStartRequest struct {
	PSBT       string `json:"psbt" validate:"required"`
	AnchorPSBT string `json:"anchor_psbt" validate:"required"`
}

type BuyCompleteRequest struct {
	PSBT            string `json:"psbt" validate:"required"`
	AnchorPSBT      string `json:"anchor_psbt" validate:"required"`
	SighashHex      string `json:"sighash_hex" validate:"required"`
	SignatureHex    string `json:"signature_hex" validate:"required"`
	AmountSatsToPay int64  `json:"amount_sats_to_pay" validate:"required"`
}

func BuyGetOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Call the service layer to fetch the buy orders.
		orders, err := GetBuyOrders()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		// Return the orders as JSON.
		return c.JSON(http.StatusOK, orders)
	}
}

// BuyStart initiates the buy process by updating the virtual PSBT.
func BuyStart(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cfg := config.GetConfig(ctx)
		pubKey := ctx.Value("public_key").(string)

		var req BuyStartRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
		}

		params := BuyStartParams{
			PSBT:         req.PSBT,
			AnchorPSBT:   req.AnchorPSBT,
			PubKey:       pubKey,
			TapdHost:     cfg.TapdHost,
			TapdMacaroon: cfg.TapdMacaroon,
		}

		response, err := StartBuyService(params, tapdClient)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, response)
	}
}

// BuyComplete finalizes the buy process by signing and finalizing the PSBT.
func BuyComplete(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cfg := config.GetConfig(ctx)

		var req BuyCompleteRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
		}

		params := BuyCompleteParams{
			PSBT:            req.PSBT,
			AnchorPSBT:      req.AnchorPSBT,
			SighashHex:      req.SighashHex,
			SignatureHex:    req.SignatureHex,
			AmountSatsToPay: req.AmountSatsToPay,
			TapdHost:        cfg.TapdHost,
			TapdMacaroon:    cfg.TapdMacaroon,
		}

		response, err := CompleteBuyService(params, tapdClient)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, response)
	}
}
