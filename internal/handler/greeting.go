package handler

import (
	"fmt"
	"net/http"
	"os"
	"rlgino/go-loki-grafana/internal/logs"

	"github.com/google/uuid"
)

type GreetingHandler struct {
	version       string
	uuidGenerated string
	logger        logs.Logger
}

func NewGreetingHandler(logger logs.Logger, uuidGenerator uuid.UUID, version string) *GreetingHandler {
	return &GreetingHandler{
		uuidGenerated: uuidGenerator.String(),
		version:       version,
		logger:        logger,
	}
}

func (handler *GreetingHandler) GetURI() string {
	return fmt.Sprintf("/%s/greeting", handler.version)
}

func (handler *GreetingHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	name := os.Getenv("NAME")
	if len(name) == 0 {
		name = "mundo"
	}
	responseBody := fmt.Sprintf("Hola %s, estas en el nodo %s", name, handler.uuidGenerated)
	_, err := writer.Write([]byte(responseBody))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	handler.logger.Log("Saludando a "+name+" desde nodo: "+handler.uuidGenerated, logs.INFO)
	writer.WriteHeader(http.StatusOK)
}
