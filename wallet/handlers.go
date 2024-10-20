package wallet

import (
	"net/http"
	"tajfi-server/config"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ConnectRequest struct {
	PublicKey string `json:"public_key" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}

// ConnectWallet handles wallet connection by validating the public key and signature.
func ConnectWallet(c echo.Context) error {
	// Extract the config from context
	cfg := config.GetConfig(c.Request().Context())

	req := new(ConnectRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	if !isValidSignature(req.PublicKey, req.Signature) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid signature",
		})
	}

	claims := jwt.MapClaims{
		"public_key": req.PublicKey,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": signedToken,
	})
}

func GetWallet(c echo.Context) error {
	wallet := Wallet{Address: GenerateNewAddress(), Balance: GetBalance()}
	return c.JSON(http.StatusOK, wallet)
}

// isValidSignature simulates signature verification (stub for demonstration).
func isValidSignature(publicKey, signature string) bool {
	// Replace with proper cryptographic signature validation logic
	return signature == "valid_signature"
}
