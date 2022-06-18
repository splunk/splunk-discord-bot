package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/splunk/splunk-discord-bot/pkg/bot"
	"github.com/splunk/splunk-discord-bot/pkg/config"
	"github.com/splunk/splunk-discord-bot/pkg/hec"
	"github.com/splunk/splunk-discord-bot/pkg/webhook"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info().Msg("Starting bot")
	cfg, err := config.ReadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Error reading config")
	}

	hecClient, err := hec.CreateClient(cfg.HecEndpoint, cfg.HecToken, cfg.InsecureSkipVerify, cfg.HecIndex)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating hec client")
	}
	log.Info().Msg("Created HEC client")

	b := bot.NewBot(cfg.Token, func(timestamp time.Time, bytes []byte) {
		err := hecClient.SendData(timestamp, bytes)
		if err != nil {
			log.Error().Err(err).Msg("Error sending data")
		}
	})

	err = b.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("bot failed to start")
	}
	log.Info().Msg("Created bot")

	s := webhook.NewServer(cfg, b)
	go func() {
		err := s.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("webhook server crashed")
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	stop := make(chan bool)
	go func() {
		<-c
		_ = s.Stop(context.Background())
		_ = b.Stop()
		_ = hecClient.Stop()
		stop <- true
	}()
	log.Info().Msg("Started bot")
	<-stop
}
