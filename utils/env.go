package utils

import (
	"cw/logger"
	"os"

	"github.com/joho/godotenv"
)

func SetENV() {
	if err := godotenv.Load(); err != nil {
		logger.GlobalLogger.Warnf("Not found ENV. Set default params(production)")
		os.Setenv("ENV", "production")
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
		os.Setenv("ENV", "development")
	}
}
