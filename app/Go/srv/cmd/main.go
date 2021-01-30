package main

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"sports/backend/srv/api/dashboard"
	"sports/backend/srv/api/server"
	"sports/backend/srv/cmd/config"
)

var cfg config.Config

// initLogger initializes the zap logger with reasonable
// defaults and replaces the global logger.
func initLogger() error {
	// Initialize the logs encoder.
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder.EncodeDuration = zapcore.StringDurationEncoder

	// Initialize the logger.
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoder,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	if err != nil {
		return err
	}

	// Then replace the globals.
	zap.ReplaceGlobals(logger)

	return nil
}

// Load up configuration.
func loadConfiguration() error {
	viper.AddConfigPath("./srv/cmd/config")
	viper.SetConfigName("configuration")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Global logging synchronizer.
	// This ensures the logged data is flushed out of the buffer before program exits.
	defer zap.S().Sync()

	err := initLogger()
	if err != nil {
		zap.S().Fatal(err)
	}

	err = loadConfiguration()
	if err != nil {
		zap.S().Fatal(err)
	}

	srv := server.Server{}
	srv.Addr = cfg.APIAddress

	run(&srv)
}

func run(srv *server.Server) error {
	c := &dashboard.Dashboard{
		ConnHub: make(map[string]*dashboard.Connection),
		Results: make(chan *dashboard.Result),
		Join:    make(chan *dashboard.Connection),
		Leave:   make(chan *dashboard.Connection),
	}

	// Disable cert verification to use self-signed certificates for internal service needs.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/dashboard", c.Handler)

	go c.Run()

	log.Printf("Server listening on %s", srv.Addr)

	zap.S().Fatal(http.ListenAndServeTLS(srv.Addr,
		"./srv/rsa.crt",
		"./srv/rsa.key",
		nil))

	return nil
}
