package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/splunk/splunk-discord-bot/pkg/bot"
	"github.com/splunk/splunk-discord-bot/pkg/config"
	"github.com/splunk/splunk-discord-bot/pkg/hec"
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

	b := bot.Bot{
		Token: cfg.Token,
		Sink: func(timestamp time.Time, bytes []byte) {
			err := hecClient.SendData(timestamp, bytes)
			if err != nil {
				log.Error().Err(err).Msg("Error sending data")
			}
		},
	}

	err = b.Start()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.Info().Msg("Created bot")


	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	stop := make(chan bool)
	go func() {
		<-c
		_ = b.Stop()
		_ = hecClient.Stop()
		stop <- true
	}()
	log.Info().Msg("Started bot")
	<-stop
}
