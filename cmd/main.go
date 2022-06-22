package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/splunk/splunk-discord-bot/pkg/bot"
	"github.com/splunk/splunk-discord-bot/pkg/config"
	"github.com/splunk/splunk-discord-bot/pkg/hec"
	"github.com/splunk/splunk-discord-bot/pkg/webhook"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal().Err(err)
	}
	logger.Info("Starting bot")
	cfg, err := config.ReadConfig()

	if err != nil {
		logger.Fatal("Error reading config", zap.Error(err))
	}

	hecClient, err := hec.CreateClient(cfg.HecEndpoint, cfg.HecToken, cfg.InsecureSkipVerify, cfg.HecIndex, logger)
	if err != nil {
		logger.Fatal("Error creating hec client", zap.Error(err))
	}
	logger.Info("Created HEC client")

	b := bot.NewBot(cfg.Token, func(timestamp time.Time, bytes []byte) {
		err := hecClient.SendData(timestamp, bytes)
		if err != nil {
			logger.Error("Error sending data", zap.Error(err))
		}
	})

	err = b.Start()

	if err != nil {
		logger.Fatal("bot failed to start", zap.Error(err))
	}
	logger.Info("Created bot")

	s := webhook.NewServer(cfg, logger, b)
	go func() {
		err := s.Start()
		if err != nil {
			logger.Fatal("webhook server crashed", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	stop := make(chan bool)
	go func() {
		<-c
		_ = s.Stop(context.Background())
		_ = b.Stop()
		_ = hecClient.Stop()
		stop <- true
	}()
	logger.Info("Started bot")
	<-stop
}
