package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	name   = "polaris"
	Output = "POLARIS_LOG"
)

var Logger *zap.Logger

func init() {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	output := os.Getenv(Output)
	if output != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, output)
	}

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Logger = log.Named(name)
}
