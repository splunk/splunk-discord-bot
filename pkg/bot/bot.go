package bot // github.com/splunk/splunk-discord-bot

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"time"
)

type Bot struct {
	Token string
	Sink  func(timestamp time.Time, data []byte)
	goBot *discordgo.Session
}

func (b *Bot) Start() error {
	log.Info().Msg("Starting bot")

	goBot, err := discordgo.New("Bot " + b.Token)

	if err != nil {
		return err
	}

	goBot.AddHandler(b.messageHandler)

	err = goBot.Open()
	if err != nil {
		return err
	}

	log.Info().Msg("Started bot")
	b.goBot = goBot

	return nil
}

func (b *Bot) Stop() error {
	return b.goBot.Close()
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling message")
		return
	}
	b.Sink(m.Timestamp, data)
}