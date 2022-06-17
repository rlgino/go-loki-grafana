package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Logger interface {
	Log(info string, level LevelLog)
}

type grafanaLogger struct {
	grafanaURL string
}

type LevelLog string

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
)

func NewLogger(grafanaURL string) Logger {
	url := fmt.Sprintf("%s/loki/api/v1/push", grafanaURL)
	return &grafanaLogger{grafanaURL: url}
}

type Message struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Stream StreamItem `json:"stream"`
	Values [][]string `json:"values"`
}

type StreamItem struct {
	Cluster  string `json:"cluster"`
	Instance string `json:"instance"`
}

func (logger *grafanaLogger) Log(info string, level LevelLog) {
	nanosecond := time.Now().UnixNano()
	bodyMsg := Message{
		Streams: []Stream{
			{
				Stream: StreamItem{
					Cluster:  "cluster-01",
					Instance: "instance-01",
				},
				Values: [][]string{
					{
						fmt.Sprintf("%d", nanosecond), info,
					},
				},
			},
		},
	}
	json_data, err := json.Marshal(bodyMsg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json_data))

	_, err = http.Post(logger.grafanaURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Printf("[Error to send to server: %v][%s:%s]\n", err.Error(), level, info)
	}
}
