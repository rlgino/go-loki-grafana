package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"rlgino/go-loki-grafana/internal/handler"
	"rlgino/go-loki-grafana/internal/logs"
)

func main() {
	logger := logs.NewLogger("http://localhost:3100")

	handlerV1 := handler.NewGreetingHandler(logger, uuid.New(), "v1")
	handlerV2 := handler.NewGreetingHandler(logger, uuid.New(), "v2")

	http.HandleFunc(handlerV1.GetURI(), handlerV1.Handle)
	http.HandleFunc(handlerV2.GetURI(), handlerV2.Handle)

	log.Println("HTTP Server running")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
