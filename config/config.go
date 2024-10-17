package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LNDHost      string `form:"LNDHost"`
	TapdHost     string `form:"TapdHost"`
	LNDMacaroon  string `form:"Macaroon"`
	TapdMacaroon string `form:"TapdMacaroon"`

	JWTSecret string `form:"JWTSecret"`
}

func (configs Config) GetConfigMap() (configMap map[string]string) {
	inrec, _ := json.Marshal(configs)
	json.Unmarshal(inrec, &configMap)
	return configMap
}

func GetConfig(ctx context.Context) (configs *Config) {
	return ctx.Value("configs").(*Config)
}

func LoadConfig(ctx context.Context) (context.Context, error) {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, falling back to .env.sample: %v", err)
		if fatalErr := godotenv.Load("env/.env.sample"); fatalErr != nil {
			log.Fatalf(fatalErr.Error())
		}
	}

	configs := &Config{
		LNDHost:      os.Getenv("LNDHost"),
		TapdHost:     os.Getenv("TapdHost"),
		LNDMacaroon:  os.Getenv("Macaroon"),
		TapdMacaroon: os.Getenv("TapdMacaroon"),
		JWTSecret:    os.Getenv("JWTSecret"),
	}

	ctx = context.WithValue(ctx, "configs", configs)

	return ctx, err
}
