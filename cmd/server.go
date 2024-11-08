package main

import (
	"context"
	"log"
	"tajfi-server/config"
	"tajfi-server/interfaces"
	"tajfi-server/wallet"
	"tajfi-server/wallet/tapd"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration into context
	ctx, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	e := echo.New()

	// Enable CORS for all domains
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Serve static assets (like CSS and JS) from /docs/swagger-ui
	e.Static("/docs", "docs")

	// Serve index.html with the correct Content-Type header
	e.GET("/docs", func(c echo.Context) error {
		return c.File("docs/index.html")
	})

	// Initialize dependencies
	// Tapd client
	tapdClient := tapd.NewTapdClient(interfaces.NewInsecureHttpClient())

	// Register wallet routes
	wallet.RegisterWalletRoutes(e, config.GetConfig(ctx), tapdClient)

	// Start the server
	if err := e.Start(":18881"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
