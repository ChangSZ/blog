package main

import (
	"fmt"

	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/router"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	conf.DefaultInit()

	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	r := router.RoutersInit()

	var opts = []http.ServerOption{ // 这里的ServerOption很多只适用于grpc protobuf
		http.Address(":8081"),
		http.Filter(handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "UPDATE"}),
			handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding",
				"X-CSRF-Token", "Authorization", "X-Auth-Token", "X-Auth-UUID", "X-Auth-Openid",
				"referrer", "Authorization", "x-client-id", "x-client-version", "x-client-type"}),
			handlers.AllowCredentials(),
			handlers.ExposedHeaders([]string{"Content-Length"}),
		)),
	}

	httpSrv := http.NewServer(opts...)
	httpSrv.HandlePrefix("/", r)

	app := kratos.New(kratos.Server(httpSrv))
	fmt.Println("开始运行")
	if err := app.Run(); err != nil {
		panic(err)
	}
}
