package functions

import (
	"github.com/stretchr/testify/assert"
	"hotalert/alert"
	"hotalert/task"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScrapeWebTask(t *testing.T) {
	var tests = []struct {
		TestName      string
		Task          task.Task
		StartServer   bool
		ServerFunc    http.HandlerFunc
		ExpectedError bool
	}{
		{
			"ServerNotOnline",
			task.Task{
				Options: task.Options{
					"keywords": []string{"keyword"},
				},
				Timeout:  10 * time.Second,
				Alerter:  alert.NewDummyAlerter(),
				Callback: nil,
			},
			false,
			func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("ok"))
			},
			true,
		},
		{
			"TestOk",
			task.Task{
				Options: task.Options{
					"keywords": []any{"keyword"},
				},
				Timeout:  10 * time.Second,
				Alerter:  alert.NewDummyAlerter(),
				Callback: nil,
			},
			true,
			func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("ok"))
			},
			false,
		},
		{
			"TestInvalidKeywords",
			task.Task{
				Options: task.Options{
					"keywords": nil,
				},
				Timeout:  10 * time.Second,
				Alerter:  alert.NewDummyAlerter(),
				Callback: nil,
			},
			true,
			func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("ok"))
			},
			true,
		},
		{
			"TestBadResponse",
			task.Task{
				Options: task.Options{
					"keywords": []string{"keyword"},
				},
				Timeout:  10 * time.Second,
				Alerter:  alert.NewDummyAlerter(),
				Callback: nil,
			},
			true,
			func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(500)
				_, _ = writer.Write([]byte("nok"))
			},
			true,
		},
		{
			"TestTimeout",
			task.Task{
				Options: task.Options{
					"keywords": []string{"keyword"},
				},
				Timeout:  0,
				Alerter:  alert.NewDummyAlerter(),
				Callback: nil,
			},
			true,
			func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("ok"))
			},
			true,
		},
	}

	for _, tv := range tests {
		t.Run(tv.TestName, func(t *testing.T) {
			testHttpServer := httptest.NewServer(tv.ServerFunc)
			if tv.StartServer {
				defer testHttpServer.Close()
			} else {
				testHttpServer.Close()
			}

			tv.Task.Options["url"] = testHttpServer.URL

			err := WebScrapeTask(&tv.Task)
			if tv.ExpectedError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
