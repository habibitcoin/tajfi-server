package wallet

import (
	"tajfi-server/config"
	"tajfi-server/middleware"

	echo "github.com/labstack/echo/v4"
)

func RegisterWalletRoutes(e *echo.Echo, cfg *config.Config) {
	api := e.Group("/api/v1")

	// No authentication for /wallet/connect
	api.POST("/wallet/connect", ConnectWallet)

	// Use auth middleware
	walletGroup := api.Group("/wallet")
	walletGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	walletGroup.GET("", GetWallet)
	walletGroup.POST("/send/decode", DecodeAddress)
	walletGroup.POST("/send/start", SendStart)
	walletGroup.POST("/send/complete", SendComplete)
	//walletGroup.GET("/balances", GetBalances)
	//walletGroup.GET("/transactions", GetTransactions)
	//walletGroup.GET("/transaction/:id", GetTransaction)
	walletGroup.POST("/receive", ReceiveAsset) // Generate an invoice to receive an asset
}
