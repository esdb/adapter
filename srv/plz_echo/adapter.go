package plz_echo

import (
	"github.com/labstack/echo"
	"reflect"
	"github.com/v2pro/plz"
	"fmt"
	"net/http"
	"github.com/v2pro/plz/srv"
	"github.com/v2pro/plz/logging"
	"errors"
)

var routerLogger = plz.LoggerOf("logger", "router")
var requestLogger = plz.LoggerOf("logger", "request")

func init() {
	srv.Executors = append(srv.Executors, StartServer)
}

func StartServer(server *srv.Server) (srv.Notifier, error) {
	e := echo.New()
	err := registerHandlers(e, "", server)
	if err != nil {
		return nil, err
	}
	httpAddress, isHttpAddressDefined := server.Properties["http_address"].(string)
	if !isHttpAddressDefined {
		return nil, errors.New("http_address not defined")
	}
	notifier := srv.NewNotifier()
	plz.Go(func() {
		err := e.Start(httpAddress)
		routerLogger.Error("server quit", "err", err)
		notifier.Stop()
	})
	plz.Go(func() {
		notifier.Wait()
		e.Listener.Close()
	})
	return notifier, nil
}

func registerHandlers(e *echo.Echo, prefix string, server *srv.Server) error {
	for _, methodProps := range server.Methods {
		name := methodProps["name"].(string)
		url := prefix + "/" + name
		decode := getDecode(methodProps)
		if decode == nil {
			return fmt.Errorf("%s: missing http_decode or echo_decode", url)
		}
		encode := getEncode(methodProps)
		if encode == nil {
			return fmt.Errorf("%s: missing http_encode or echo_encode", url)
		}
		handle := getHandle(methodProps)
		if handle == nil {
			return fmt.Errorf("%s: missing handle", url)
		}
		httpMethod, ok := methodProps["method"].(string)
		if !ok {
			httpMethod = "GET"
		}
		routerLogger.Info("register http route", "url", url)
		switch httpMethod {
		case "GET":
			e.GET(url, func(ctx echo.Context) error {
				if requestLogger.ShouldLog(logging.LEVEL_DEBUG) {
					requestLogger.Debug("handle request", "url", ctx.Request().URL.String())
				}
				request, err := decode(ctx)
				if err != nil {
					return err
				}
				httpContext := ctx.Request().Context()
				ret := handle.Call([]reflect.Value{reflect.ValueOf(httpContext), reflect.ValueOf(request)})
				resp := ret[0].Interface()
				err, _ = ret[1].Interface().(error)
				return encode(ctx, resp, err)
			})
		default:
			return fmt.Errorf("unknown http method %v", methodProps["method"])
		}
	}
	for _, subServer := range server.SubServers {
		err := registerHandlers(e, prefix+"/"+subServer.Properties["name"].(string), subServer)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDecode(methodProps map[string]interface{}) func(echo.Context) (interface{}, error) {
	decode, found := methodProps["echo_decode"].(func(echo.Context) (interface{}, error))
	if found {
		return decode
	}
	httpDecode, found := methodProps["http_decode"].(func(*http.Request) (interface{}, error))
	if found {
		return func(ctx echo.Context) (interface{}, error) {
			return httpDecode(ctx.Request())
		}
	}
	requestType := reflect.TypeOf(methodProps["handle"]).In(1)
	return func(ctx echo.Context) (interface{}, error) {
		// default to use generic copy to bind request
		req := reflect.New(requestType).Interface()
		err := plz.Copy(req, ctx)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
}

func getEncode(methodProps map[string]interface{}) func(echo.Context, interface{}, error) error {
	encode, found := methodProps["echo_encode"].(func(echo.Context, interface{}, error) error)
	if found {
		return encode
	}
	httpEncode, found := methodProps["http_encode"].(func(http.ResponseWriter, interface{}, error) error)
	if found {
		return func(ctx echo.Context, resp interface{}, err error) error {
			return httpEncode(ctx.Response().Writer, resp, err)
		}
	}
	return func(ctx echo.Context, resp interface{}, respErr error) error {
		// default to use generic copy to bind request
		if respErr != nil {
			return plz.Copy(ctx, respErr)
		}
		return plz.Copy(ctx, resp)
	}
}

func getHandle(methodProps map[string]interface{}) *reflect.Value {
	handle, _ := methodProps["handle"]
	if handle == nil {
		return nil
	}
	val := reflect.ValueOf(handle)
	return &val
}
