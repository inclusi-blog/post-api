package main

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/tracing"
	"net/http"
	. "post-api/init"
)

func main() {
	configData := LoadConfig()
	router := CreateRouter(configData)
	tracing.Init(configData.TracingServiceName, configData.TracingOCAgentHost)
	err := http.ListenAndServe(":8080", tracing.WithTracing(router, "/api/post/healthz"))
	if err != nil {
		logging.GetLogger(context.TODO()).Error("Could not start the server", err)
	}
}
