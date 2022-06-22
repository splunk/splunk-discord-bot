package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/splunk/splunk-discord-bot/pkg/bot"
	"github.com/splunk/splunk-discord-bot/pkg/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	httpAddr           = ":8080"
	serverReadTimeout  = time.Second * 10
	serverWriteTimeout = time.Second * 60
	serverIdleTimeout  = time.Second * 120
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

func NewServer(cfg *config.Config, logger *zap.Logger, bot bot.Bot) Server {
	return &serverImpl{
		bot:        bot,
		webhooks:   cfg.WebHooks,
		listenAddr: cfg.WebhookListenAddr,
		logger:     logger,
	}
}

type serverImpl struct {
	bot        bot.Bot
	webhooks   []*config.WebhookConfig
	webServer  *http.Server
	listenAddr string
	logger     *zap.Logger
}

func (s *serverImpl) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(405)
		s.logger.Debug("invalid verb", zap.String("remoteAddr", req.RemoteAddr), zap.String("method", req.Method))
		return
	}
	wh := req.URL.Query().Get("webhook")
	var cfg *config.WebhookConfig
	for _, whc := range s.webhooks {
		if whc.ID == wh {
			cfg = whc
			break
		}
	}
	if cfg == nil {
		resp.WriteHeader(404)
		s.logger.Debug("no webhook defined", zap.String("remoteAddr", req.RemoteAddr), zap.String("webhook", wh))
		return
	}

	var alertRequest AlertRequest
	err := json.NewDecoder(req.Body).Decode(&alertRequest)
	if err != nil {
		s.logger.Debug("bad request", zap.String("remoteAddr", req.RemoteAddr), zap.Error(err))
		resp.WriteHeader(400)
		return
	}
	err = s.bot.SendMessage(cfg.Channel, fmt.Sprintf("%s - see results %s", alertRequest.SearchName, alertRequest.ResultsLink))
	if err != nil {
		s.logger.Error("error sending message to Discord", zap.String("remoteAddr", req.RemoteAddr), zap.Error(err))

		resp.WriteHeader(500)
		return
	}
	resp.WriteHeader(200)
}

func (s *serverImpl) Start() error {
	listen := s.listenAddr
	if listen == "" {
		listen = httpAddr
	}
	s.webServer = &http.Server{
		Addr:         listen,
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
