package util

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func SetRequest(context *gin.Context, request interface{}) {
	requestBody, _ := Encode(request)
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))
}

func Encode(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseResponse(recorder *httptest.ResponseRecorder, response interface{}) {
	bodyBytes, _ := ioutil.ReadAll(recorder.Body)
	json.Unmarshal(bodyBytes, &response)
}
