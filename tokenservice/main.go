package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stillfox-lee/learn-microservice/token/endpoints"
	"github.com/stillfox-lee/learn-microservice/token/infra"
	"github.com/stillfox-lee/learn-microservice/token/model"
	"github.com/stillfox-lee/learn-microservice/token/services"
	"github.com/stillfox-lee/learn-microservice/token/transports"
)

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	// init flag
	var (
		redisHost     = flag.String("redis-host", "localhost", "redis host")
		redisPort     = flag.String("redis-port", "6379", "redis port")
		redisPassword = flag.String("redis-password", "", "redis password")
		redisDB       = flag.Int("redis-db", 0, "redis db")
		serverHost    = flag.String("server-host", "0.0.0.0", "server host")
		serverPort    = flag.String("server-port", "80", "server port")
		serviceName   = flag.String("service-name", "token-service", "service name")
		consulHost    = flag.String("consul-host", "localhost", "consul host")
		consulPort    = flag.Int("consul-port", 8500, "consul port")
	)
	flag.Parse()

	// init logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dClient, err := infra.NewDiscoveryClient(*consulHost, *consulPort)
	if err != nil {
		logger.Log("init discovery client fail", err)
		os.Exit(-1)
	}

	// init redis client
	redisClient, err := infra.InitRedisClient(*redisHost, *redisPort, *redisPassword, *redisDB)
	if err != nil {
		logger.Log("init redis client fail", err)
		os.Exit(-1)
	}
	dao := model.MakeDao(redisClient)
	svc := services.MakeTokenService(dao, logger)

	e := endpoints.TokenEndpoints{
		CreateTokenEndpoint: endpoints.MakeCreateTokenEndpoint(svc),
		GetTokenEndpoint:    endpoints.MakeGetTokenEndpoint(svc),
		HealthEndpoint:      endpoints.MakeHealthEndpoint(svc),
	}

	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(transports.EncodeError),
	}

	r.Methods("POST").Path("/token/create").Handler(kithttp.NewServer(
		e.CreateTokenEndpoint,
		transports.DecodeCreateRequest,
		transports.EncodeJSONResponse,
		options...,
	))
	r.Methods("GET").Path("/token/get").Handler(kithttp.NewServer(
		e.GetTokenEndpoint,
		transports.DecodeGetRequest,
		transports.EncodeJSONResponse,
		options...,
	))
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		e.HealthEndpoint,
		transports.DecodeHealthRequest,
		transports.EncodeJSONResponse,
		options...,
	))

	addr := fmt.Sprintf("%s:%s", *serverHost, *serverPort)

	exitCh := make(chan error)
	go func() {
		logger.Log("server", *serviceName, addr)
		err := http.ListenAndServe(addr, r)
		if err != nil {
			exitCh <- err
		}
	}()
	// register service to consul
	instanceId := fmt.Sprintf("%s-%s", *serviceName, uuid.New().String())
	dClient.Register(*serviceName, instanceId, getIP(), *serverPort, "/health", nil, nil)

	// handle exit signal
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		exitCh <- fmt.Errorf("%s", <-sigCh)
	}()

	// wait for exit
	err = <-exitCh
	dClient.Deregister()
	logger.Log("exit", err)

}
