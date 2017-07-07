package main

import (
	"github.com/v2pro/plz/logging"
	"github.com/v2pro/plz/srv"
	"context"
	_ "github.com/v2pro/plz_echo"
	_ "github.com/v2pro/plz/lang/nativeacc"
)

func main() {
	logging.Providers = append(logging.Providers, func(loggerKv []interface{}) logging.Logger {
		return logging.NewStderrLogger(loggerKv, logging.LEVEL_DEBUG)
	})
	signal, err := srv.BuildServer("http_address", "127.0.0.1:9000").
		Method("example", handleMyRequest).
		Start()
	if err != nil {
		panic(err.Error())
	}
	signal.Wait()
}

type MyRequest struct {
	Field string
}

type MyResponse struct {
	Field string
}

func handleMyRequest(ctx context.Context, req MyRequest) (MyResponse, error) {
	return MyResponse{req.Field}, nil
}
