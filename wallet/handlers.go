package wallet

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tajfi-server/config"
	"tajfi-server/wallet/tapd"
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

func GetBalances(c echo.Context) error {
	// Extract config from context
	ctx := c.Request().Context()
	cfg := config.GetConfig(ctx)
	pubKey := ctx.Value("public_key").(string)

	utxos, err := tapd.GetUtxos(cfg.TapdHost, cfg.TapdMacaroon)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch balances from tapd: "+err.Error())
	}

	balances, err := constructWalletBalancesResponse(utxos, pubKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to construct wallet balances: "+err.Error())
	}
	return c.JSON(http.StatusOK, balances)
}

// ConstructWalletBalancesResponse constructs a WalletBalancesResponse by finding
// every Asset with a ScriptKey matching the provided value, summing up their amounts,
// and grouping them by AssetGenesis.
func constructWalletBalancesResponse(utxos *tapd.GetUtxosResponse, scriptKey string) (*tapd.WalletBalancesResponse, error) {
	assetBalances := make(map[string]tapd.AssetBalance)

	for _, utxo := range utxos.ManagedUtxos {
		for _, asset := range utxo.Assets {
			genesisID := asset.AssetGenesis.AssetID
			amount := 0

			if asset.ScriptKey == ("02" + scriptKey) {

				log.Println(asset)
				var err error
				amount, err = strconv.Atoi(asset.Amount)
				if err != nil {
					return nil, fmt.Errorf("invalid amount: %v", err)
				}
			}

			if balance, exists := assetBalances[genesisID]; exists {
				existingAmount, err := strconv.Atoi(balance.Balance)
				if err != nil {
					return nil, fmt.Errorf("invalid existing balance: %v", err)
				}
				balance.Balance = strconv.Itoa(existingAmount + amount)
				assetBalances[genesisID] = balance
			} else {
				assetBalances[genesisID] = tapd.AssetBalance{
					AssetGenesis: asset.AssetGenesis,
					Balance:      strconv.Itoa(amount),
				}
			}
		}
	}

	return &tapd.WalletBalancesResponse{AssetBalances: assetBalances}, nil
}

func GetTransfers(c echo.Context) error {
	// Extract config from context
	cfg := config.GetConfig(c.Request().Context())

	balances, err := tapd.GetTransfers(cfg.TapdHost, cfg.TapdMacaroon)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch balances from tapd: "+err.Error())
	}
	return c.JSON(http.StatusOK, balances)
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
