package logger

import "go.uber.org/zap"

// NewLogger initializes a new zap logger with production configuration.
// It can be customized to set different log levels or formats as needed.
// Currently, it returns a production logger with default settings.
// Uncomment the configuration lines to customize the logger further.
func NewLogger() (*zap.Logger, error) {
	// cfg := zap.NewProductionConfig()
	// cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	// return cfg.Build()

	return zap.NewProduction()
}
