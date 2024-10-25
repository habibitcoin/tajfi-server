package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	LNDHost      string `form:"LNDHost"`
	TapdHost     string `form:"TapdHost"`
	LNDMacaroon  string `form:"LNDMacaroon"`
	TapdMacaroon string `form:"TapdMacaroon"`

	JWTSecret string `form:"JWTSecret"`

	TaprootSigsDir string `form:"TaprootSigsDir"`

	DemoMode         bool   `form:"DemoMode"`
	DemoAmount       int    `form:"DemoAmount"`
	DemoTapdHost     string `form:"DemoTapdHost"`
	DemoTapdMacaroon string `form:"DemoTapdMacaroon"`
}

func (configs Config) GetConfigMap() (configMap map[string]string) {
	inrec, _ := json.Marshal(configs)
	json.Unmarshal(inrec, &configMap)
	return configMap
}

func GetConfig(ctx context.Context) (configs *Config) {
	ctx, _ = LoadConfig(ctx)
	return ctx.Value("configs").(*Config)
}

func LoadConfig(ctx context.Context) (context.Context, error) {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, falling back to .env.sample: %v", err)
		if fatalErr := godotenv.Load(".env.sample"); fatalErr != nil {
			log.Fatalf(fatalErr.Error())
		}
	}

	demoModeStr := os.Getenv("DemoMode")
	demoMode := false
	if demoModeStr == "true" {
		demoMode = true
	}

	demoAmountStr := os.Getenv("DemoAmount")
	demoAmount := 0
	if demoAmountStr != "" {
		// convert string to int
		demoAmount, err = strconv.Atoi(demoAmountStr)
		if err != nil {
			demoMode = false
		}
	}

	configs := &Config{
		LNDHost:          os.Getenv("LNDHost"),
		TapdHost:         os.Getenv("TapdHost"),
		LNDMacaroon:      os.Getenv("LNDMacaroon"),
		TapdMacaroon:     os.Getenv("TapdMacaroon"),
		JWTSecret:        os.Getenv("JWTSecret"),
		TaprootSigsDir:   os.Getenv("TaprootSigsDir"),
		DemoMode:         demoMode,
		DemoAmount:       demoAmount,
		DemoTapdHost:     os.Getenv("DemoTapdHost"),
		DemoTapdMacaroon: os.Getenv("DemoTapdMacaroon"),
	}

	ctx = context.WithValue(ctx, "configs", configs)

	return ctx, err
}
