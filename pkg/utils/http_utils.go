package utils

import (
	"e-klinik/pkg"
	"e-klinik/pkg/logging"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-zoox/fetch"
)

func GetFileRequest(url string) (*http.Response, error) {
	// Step 1: Fetch the file from the URL
	response, err := http.Get(url)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "Failed to fetch file")
	}

	// Check if the response status is 200 OK
	if response.StatusCode != http.StatusOK {
		return nil, pkg.NewErrorf(pkg.ErrorCodeUnknown, fmt.Sprintf("Failed to fetch image, status code: %d", response.StatusCode))

	}
	return response, nil
}

func GetHttpRequest[T any](url string, log logging.Logger) (T, error) {
	var zeroValue T
	response, err := fetch.Get(url, &fetch.Config{
		Headers: fetch.Headers{"Accept": "*/*",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
			"User-Agent":      "Go-http-client/1.1"},
	})

	if err != nil {
		log.Error(logging.Http, logging.HttpError, err.Error(), nil)
		return zeroValue, err
	}

	if !response.Ok() {
		log.Error(logging.Http, logging.HttpError, fmt.Sprintf("HTTP request failed with status code %d: %s", response.StatusCode(), response.Body), nil)
		return zeroValue, fmt.Errorf("HTTP request failed with status code %d: %s", response.StatusCode(), response.Body)
	}
	var res T
	err = json.Unmarshal(response.Body, &res)
	if err != nil {
		log.Error(logging.Http, logging.HttpError, err.Error(), nil)
		return zeroValue, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "unknown error")
	}

	return res, nil
}
