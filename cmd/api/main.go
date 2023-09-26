package main

import (
	"context"
	"external-metrics/config"
	coingeckoprovider "external-metrics/pkg/coingecko"
	"external-metrics/pkg/tools/logging"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	configPath := flag.String("c", "", "config/config.yaml")
	flag.Parse()
	if *configPath == "" {
		log.Fatalf("config path flag -c is required")
	}
	// config set up
	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("Can't get config from yaml: %v", err)
	}

	f, err := os.OpenFile(cfg.App.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Can't open log file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print("can't close log file: ", err)
		}
	}()
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.Println("start blockchain-statistic-metrics")

	logLevel, err := logrus.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		log.Panic(err)
	}

	logger := logging.InitLoggerNew(logLevel, io.MultiWriter(os.Stdout, f))
	logger.Debug("Application is setting up")

	// coingecko service
	coingeckoProvider, err := coingeckoprovider.NewProvider(cfg.Coingecko.APIAddress, logger)
	if err != nil {
		log.Panic("Create coingeckoProvider error: ", err)
	}

	// bootstrap server
	workspaceAPI := bootstrapAPI(coingeckoProvider, cfg, logger)

	// graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// start server
	logger.Infof("Starting blockchain statistic proxy API server at %s...", cfg.App.APIAddress)
	workspaceAPIDone := startServer(workspaceAPI)

	select {
	case <-signalChan:
		logger.Infof("Received system interrupt, graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err = workspaceAPI.Shutdown(ctx)
		if err != nil {
			logger.Errorf("Error on graceful shutdown: %v", err)

		}
	case err = <-workspaceAPIDone:
		logger.Errorf("blockchain statistic proxy API stopped: %v", err)
	}
}
