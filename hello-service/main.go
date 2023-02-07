package main

import (
	"log"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	svc := &helloService{}
	ep := MakeHelloEndpoint(svc)

	route := mux.NewRouter()
	route.Methods("Get").Path("/hello").Handler(kithttp.NewServer(
		ep,
		decodeHelloRequest,
		encodeResponse,
	))

	log.Fatal(http.ListenAndServe(":8080", route))
}
