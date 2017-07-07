package plz_http

import (
	"testing"
	"github.com/v2pro/plz"
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
)

func Test_http_get(t *testing.T) {
	should := require.New(t)
	ExecutorProviders = append(ExecutorProviders, func(serviceName string, methodName string, kv ...interface{}) HttpExecutor {
		targetUrl, _ := url.Parse("http://www.douban.com")
		return &dummyExecutor{
			method:    "GET",
			host:      serviceName,
			targetUrl: targetUrl,
		}
	})
	client := plz.ClientOf("douban", "index")
	req := MyRequest{}
	resp := MyResponse{}
	err := client.Call(context.TODO(), req, &resp)
	should.Nil(err)
}

type MyRequest struct {
}

type MyResponse struct {
}

type dummyExecutor struct {
	method    string
	host      string
	targetUrl *url.URL
	executor  HttpExecutor
}

func (executor *dummyExecutor) Do(req *http.Request) (*http.Response, error) {
	req.Method = executor.method
	req.Host = executor.host
	req.URL = executor.targetUrl
	return executor.executor.Do(req)
}
