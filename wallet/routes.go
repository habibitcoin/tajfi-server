package wallet

import (
	"tajfi-server/config"
	"tajfi-server/middleware"
	"tajfi-server/wallet/tapd"

	echo "github.com/labstack/echo/v4"
)

func RegisterWalletRoutes(e *echo.Echo, cfg *config.Config, tapdClient tapd.TapdClientInterface) {
	api := e.Group("/api/v1")

	// No authentication for /wallet/connect
	api.POST("/wallet/connect", ConnectWallet)

	// Use auth middleware
	walletGroup := api.Group("/wallet")
	walletGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	walletGroup.GET("", GetWallet)
	walletGroup.POST("/send/decode", DecodeAddress(tapdClient))
	walletGroup.POST("/send/start", SendStart(tapdClient))
	walletGroup.POST("/send/complete", SendComplete(tapdClient))
	walletGroup.GET("/balances", GetBalances(tapdClient))
	walletGroup.GET("/transfers", GetTransfers(tapdClient))
	//walletGroup.GET("/transaction/:id", GetTransaction)
	walletGroup.POST("/receive", ReceiveAsset(tapdClient)) // Generate an invoice to receive an asset
}
