package plz_http

import (
	"github.com/v2pro/plz/clt"
	"fmt"
	"context"
	"github.com/v2pro/plz"
	"net/http"
)

type HttpExecutor interface {
	Do(req *http.Request) (*http.Response, error)
}

var ExecutorProviders = []func(serviceName string, methodName string, kv ...interface{}) HttpExecutor{
}

func init() {
	clt.Providers = append(clt.Providers, func(serviceName string, methodName string, kv ...interface{}) clt.Client {
		for _, provider := range ExecutorProviders {
			executor := provider(serviceName, methodName, kv...)
			if executor != nil {
				return &httpClientAdapter{executor}
			}
		}
		panic(fmt.Sprintf("no executor defined for %s %s", serviceName, methodName))
	})
}

type httpClientAdapter struct {
	executor HttpExecutor
}

func (adapter *httpClientAdapter) Call(ctx context.Context, req interface{}, resp interface{}) error {
	httpReq := &http.Request{}
	httpReq = httpReq.WithContext(ctx)
	err := plz.Copy(httpReq, req)
	if err != nil {
		return err
	}
	httpResp, err := adapter.executor.Do(httpReq)
	if err != nil {
		return err
	}
	err = plz.Copy(resp, httpResp)
	if err != nil {
		return err
	}
	return nil
}
