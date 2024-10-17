package main

import (
	"context"
	"log"
	"tajfi-server/config"
	"tajfi-server/wallet"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration into context
	ctx, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	e := echo.New()

	// Register wallet routes
	wallet.RegisterWalletRoutes(e, ctx)

	// Start the server
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
