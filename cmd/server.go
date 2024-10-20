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

	// Serve static assets (like CSS and JS) from /docs/swagger-ui
	e.Static("/docs", "docs")

	// Serve index.html with the correct Content-Type header
	e.GET("/docs", func(c echo.Context) error {
		return c.File("docs/index.html")
	})

	// Register wallet routes
	wallet.RegisterWalletRoutes(e, config.GetConfig(ctx))

	// Start the server
	if err := e.Start(":18881"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
