package init

import (
	"github.com/gola-glitch/gola-utils/http/request"
	"go.opencensus.io/plugin/ochttp"
	"net/http"
	"post-api/configuration"
	"time"
)

var (
	client         http.Client
	requestBuilder request.HttpRequestBuilder
)

func HttpClient(data *configuration.ConfigData) {
	client = http.Client{
		Transport: &ochttp.Transport{},
		Timeout:   time.Second * time.Duration(data.RequestTimeout),
	}
	requestBuilder = request.NewHttpRequestBuilder(&client)
}
