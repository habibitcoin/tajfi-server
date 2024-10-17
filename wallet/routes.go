package wallet

import (
	"context"
	"tajfi-server/middleware"

	echo "github.com/labstack/echo/v4"
)

func RegisterWalletRoutes(e *echo.Echo, ctx context.Context) {
	api := e.Group("/api/v1")

	// No authentication for /wallet/connect
	api.POST("/wallet/connect", func(c echo.Context) error {
		// Inject the context into the request
		c.SetRequest(c.Request().WithContext(ctx))
		return ConnectWallet(c)
	})

	// Use auth middleware
	walletGroup := api.Group("/wallet")
	walletGroup.Use(middleware.AuthMiddleware(ctx))

	walletGroup.GET("", GetWallet)
	//walletGroup.POST("/send", SendAsset)
	//walletGroup.GET("/balances", GetBalances)
	//walletGroup.GET("/transactions", GetTransactions)
	//walletGroup.GET("/transaction/:id", GetTransaction)
	//walletGroup.POST("/receive", ReceiveAsset) // Generate an invoice to receive an asset
}
