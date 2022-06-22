package hec

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateClient(t *testing.T) {
	messageChan := make(chan []byte)
	go func() {
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			messageChan <- body
		})
		http.Handle("/", handler)

		err := http.ListenAndServe(":8080", nil)
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	hecClient, err := CreateClient("http://localhost:8080/", "1111-1111", true, "index", zap.NewNop())
	assert.NoError(t, err)
	logs := plog.NewLogs()
	logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty().Body().SetStringVal("hello")
	err = hecClient.SendLogs(logs)
	assert.NoError(t, err)

	message := <-messageChan
	assert.Equal(t, `{"host":"unknown","index":"index","event":"hello"}`, string(message))
}
