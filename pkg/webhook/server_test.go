package webhook

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/splunk/splunk-discord-bot/pkg/config"
	"github.com/splunk/splunk-discord-bot/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const correctJson = `{
	"result": {
		"sourcetype" : "mongod",
		"count" : "8"
	},
	"sid" : "scheduler_admin_search_W2_at_14232356_132",
	"results_link" : "http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132",
	"search_name" : "alert_name",
	"owner" : "admin",
	"app" : "search"
}`

func Test_Server_ServeHTTP(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name       string
		req        *http.Request
		statusCode int
		botFn      func(b *mocks.MockBot)
	}{
		{
			name:       "GET",
			req:        httptest.NewRequest("GET", "/", strings.NewReader("foo")),
			statusCode: 405,
		},
		{
			name:       "bad webhook",
			req:        httptest.NewRequest("POST", "/?webhook=bar", strings.NewReader("foo")),
			statusCode: 404,
		},
		{
			name:       "correct webhook bad json",
			req:        httptest.NewRequest("POST", "/?webhook=foo", strings.NewReader("{foo")),
			statusCode: 400,
		},
		{
			name:       "correct webhook correct json bot error",
			req:        httptest.NewRequest("POST", "/?webhook=foo", strings.NewReader(correctJson)),
			statusCode: 500,
			botFn: func(b *mocks.MockBot) {
				b.EXPECT().SendMessage("somechannel", "alert_name - see results http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132").Return(errors.New("some error"))
			},
		},
		{
			name:       "correct webhook correct json no error",
			req:        httptest.NewRequest("POST", "/?webhook=foo", strings.NewReader(correctJson)),
			statusCode: 200,
			botFn: func(b *mocks.MockBot) {
				b.EXPECT().SendMessage("somechannel", "alert_name - see results http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132").Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockBot := mocks.NewMockBot(mockCtrl)
			s := NewServer(&config.Config{
				Token:              "1234",
				HecEndpoint:        "",
				HecToken:           "",
				InsecureSkipVerify: false,
				HecIndex:           "",
				WebHooks: []*config.WebhookConfig{
					{
						ID:      "foo",
						Channel: "somechannel",
					},
				},
			}, zap.NewNop(), mockBot).(*serverImpl)
			recorder := httptest.NewRecorder()
			if test.botFn != nil {
				test.botFn(mockBot)
			}
			s.ServeHTTP(recorder, test.req)
			assert.Equal(t, test.statusCode, recorder.Code)
		})
	}
}
