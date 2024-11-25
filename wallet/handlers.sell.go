package wallet

import (
	"net/http"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
)

type SellStartRequest struct {
	AssetID      string `json:"asset_id" validate:"required"`
	AmountToSell int64  `json:"amount_to_sell" validate:"required"`
}

type SellStartResponse struct {
	tapd.FundVirtualPSBTResponse
}

type SellCompleteRequest struct {
	PSBT                string `json:"psbt" validate:"required"`
	SighashHex          string `json:"sighash_hex" validate:"required"`
	SignatureHex        string `json:"signature_hex" validate:"required"`
	AmountSatsToReceive int64  `json:"amount_sats_to_receive" validate:"required"`
}

type SellCompleteResponse struct {
	SignedVirtualPSBT  string `json:"signed_virtual_psbt"`
	ModifiedAnchorPSBT string `json:"modified_anchor_psbt"`
}

// SellStart initiates the sell process by creating a virtual PSBT.
func SellStart(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cfg := config.GetConfig(ctx)
		pubKey := ctx.Value("public_key").(string)

		var req SellStartRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
		}

		// Call the service to start the sell process.
		params := SellStartParams{
			PubKey:       pubKey,
			AssetID:      req.AssetID,
			AmountToSell: req.AmountToSell,
			TapdHost:     cfg.TapdHost,
			TapdMacaroon: cfg.TapdMacaroon,
		}

		response, err := StartSellService(params, tapdClient)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, response)
	}
}

// SellComplete finalizes the sell process by signing the virtual PSBT and preparing the anchoring template.
func SellComplete(tapdClient tapd.TapdClientInterface) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cfg := config.GetConfig(ctx)

		var req SellCompleteRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
		}

		// Call the service to complete the sell process.
		params := SellCompleteParams{
			Signature:           req.SignatureHex,
			PSBT:                req.PSBT,
			SighashHex:          req.SighashHex,
			SignatureHex:        req.SignatureHex,
			AmountSatsToReceive: req.AmountSatsToReceive,
			TapdHost:            cfg.TapdHost,
			TapdMacaroon:        cfg.TapdMacaroon,
		}

		response, err := CompleteSellService(params, tapdClient)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, response)
	}
}
