package hec

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/model/pdata"
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

		http.ListenAndServe(":8080", nil)
	}()

	hecClient, err := CreateClient("http://localhost:8080/", "1111-1111")
	assert.NoError(t, err)
	logs := pdata.NewLogs()
	logs.ResourceLogs().AppendEmpty().InstrumentationLibraryLogs().AppendEmpty().LogRecords().AppendEmpty().Body().SetStringVal("hello")
	hecClient.sendLogs(logs)

	message := <-messageChan
	assert.Equal(t, `{"host":"unknown","event":"hello"}`, string(message))
}
