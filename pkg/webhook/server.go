package webhook

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/splunk/splunk-discord-bot/pkg/bot"
	"github.com/splunk/splunk-discord-bot/pkg/config"
	"net/http"
	"time"
)

const (
	httpAddr = ":8080"
	serverReadTimeout  = time.Second * 10
	serverWriteTimeout = time.Second * 60
	serverIdleTimeout  = time.Second * 120
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

func NewServer(cfg config.Config, bot bot.Bot) Server {
	return &serverImpl {
		bot: bot,
		webhooks: cfg.WebHooks,
	}
}

type serverImpl struct {
	bot           bot.Bot
	webhooks      []config.WebhookConfig
	webServer     *http.Server
}

func (s *serverImpl) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(405)
		return
	}
	wh := req.URL.Query().Get("webhook")
	var cfg *config.WebhookConfig
	for _, whc := range s.webhooks {
		if whc.ID == wh {
			cfg = &whc
			break
		}
	}
	if cfg == nil {
		resp.WriteHeader(404)
		return
	}

	var alertRequest AlertRequest
	err := json.NewDecoder(req.Body).Decode(&alertRequest)
	if err != nil {
		log.Error().Err(err).Msg("Error reading the JSON of the request")
		resp.WriteHeader(400)
		return
	}

	err = s.bot.SendMessage(cfg.Channel, alertRequest.Result.Count)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message to Discord")
		resp.WriteHeader(500)
		return
	}
	resp.WriteHeader(200)
}

func (s *serverImpl) Start() error {
	s.webServer = &http.Server{
		Addr:         httpAddr,
		Handler:      s,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}
	err := s.webServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *serverImpl) Stop(ctx context.Context) error {
	return s.webServer.Shutdown(ctx)
}

type AlertRequest struct {
	Result struct {
		Sourcetype string `json:"sourcetype"`
		Count      string `json:"count"`
	} `json:"result"`
	Sid         string      `json:"sid"`
	ResultsLink string      `json:"results_link"`
	SearchName  interface{} `json:"search_name"`
	Owner       string      `json:"owner"`
	App         string      `json:"app"`
}