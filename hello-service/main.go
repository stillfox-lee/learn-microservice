package main

import (
	"log"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	svc := &helloService{}
	ep := MakeHelloEndpoint(svc)

	// decorate ratelimit
	ratelimit := limitMiddleware{
		timer: 5 * time.Second,
		burst: 3,
	}
	ep = ratelimit.wrap(ep)

	route := mux.NewRouter()
	route.Methods("Get").Path("/hello").Handler(kithttp.NewServer(
		ep,
		decodeHelloRequest,
		encodeResponse,
	))

	log.Fatal(http.ListenAndServe(":8080", route))
}
