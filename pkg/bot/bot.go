package bot // github.com/splunk/splunk-discord-bot

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_bot.go -package=mocks . Bot
type Bot interface {
	Start() error
	Stop() error
	SendMessage(channel string, msg string) error
}

type botImpl struct {
	token string
	sink  func(timestamp time.Time, data []byte)
	goBot *discordgo.Session
}

func NewBot(token string, sinkFn func(timestamp time.Time, bytes []byte)) Bot {
	return &botImpl{
		token: token,
		sink:  sinkFn,
	}
}

func (b *botImpl) Start() error {
	log.Info().Msg("Starting bot")

	goBot, err := discordgo.New("Bot " + b.token)

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

func (b *botImpl) Stop() error {
	return b.goBot.Close()
}

func (b *botImpl) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	data, err := json.Marshal(m)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling message")
		return
	}
	b.sink(m.Timestamp, data)
}

func (b *botImpl) SendMessage(channel string, msg string) error {
	message, err := b.goBot.ChannelMessageSend(channel, msg)
	if err != nil {
		return err
	}
	log.Info().Str("channel", channel).Time("time", message.Timestamp).Str("message", message.Content).Msg("Sending alert message")
	return nil
}
