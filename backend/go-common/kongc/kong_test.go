package kongc

import (
	"fmt"
	"go-common/log"
	"go-common/utils/httpc"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestRegister9090(t *testing.T) {
	listenServer(9090)
}

func TestRegister9091(t *testing.T) {
	listenServer(9091)
}

func TestRegister9092(t *testing.T) {
	listenServer(9092)
}

func listenServer(port int) {
	devZapConfig := log.DevZapConf()
	prodZapConf := log.ProdZapConf()
	prodZapConf.IsJson = true
	log.Init(devZapConfig, prodZapConf)
	loggingMiddleware := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			// 打印请求地址
			zap.S().Infof("服务器端口:%d请求地址: %s %s", port, request.Method, request.URL.Path)
			handler(writer, request)
		}
	}

	// 应用中间件到所有路由
	http.HandleFunc("/", loggingMiddleware(func(writer http.ResponseWriter, request *http.Request) {
		// 根据不同的路径返回不同的响应
		switch request.URL.Path {
		case "/user":
			writer.Header().Set("Content-Type", "application/json")
			_, err := writer.Write([]byte("{\"name\":\"张三\"}"))
			if err != nil {
				zap.S().Errorw("log.Println", zap.Error(err))
			}
		default:
			// 对于其他路径，可以返回200
			_, err := writer.Write([]byte("request.URL.Path"))
			if err != nil {
				zap.S().Errorw("log.Println", zap.Error(err))
			}

		}
	}))
	// 应用中间件到 /user 路由
	http.HandleFunc("/user", loggingMiddleware(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte("{\"name\":\"张三\"}"))
		if err != nil {
			zap.S().Errorw("log.Println", zap.Error(err))
		}
	}))
	client := Init("http://127.0.0.1:8001", "")
	err := client.RegisterHttpWithOption(&RegisterHttpOption{
		Host:        "host.docker.internal",
		Port:        port,
		ServiceName: "test",
		RouteOptions: []RouteOption{
			{
				Name:      "user",
				Paths:     []string{"/test", "/user"},
				Methods:   []httpc.Method{httpc.MethodGet},
				StripPath: false,
			},
		},
	})
	if err != nil {
		zap.S().Errorf("err: %+v", err)
	}
	// 启动web服务，监听9090端口
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		zap.S().Fatalw("log.Fatal", zap.Error(err))
	}
}
