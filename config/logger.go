package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger() {
	file, _ := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	Logger = zerolog.New(file).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339
}
