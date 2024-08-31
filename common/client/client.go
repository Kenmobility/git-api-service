package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

var DefaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

type Method string

// Supported HTTP verbs.
const (
	Get Method = "GET"
)

// Client is an enhanced http.Client.
// Can plug in cache, etc
type RestClient struct{}

func NewRestClient() *RestClient {
	return &RestClient{}
}

// Request holds the request to an API Call.
type Request struct {
	Method      Method
	BaseURL     string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

type Response struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}

var ErrNoRequest = errors.New("nil http.Request received")

// MakeRequest uses the DefaultClient to send the request and returns the response.
func makeRequest(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, ErrNoRequest
	}

	var (
		resp = new(http.Response)
		err  error
		url  = req.URL.String()
	)

	start := time.Now()

	resp, err = DefaultHTTPClient.Do(req)
	if err != nil {
		log.Info().Msgf("Failed to do request; url: %s method: %s, err: %v", url, req.Method, err)

		if errors.Is(err, context.DeadlineExceeded) {
			return resp, fmt.Errorf("dal: %w", err)
		}

		return resp, err
	}

	log.Info().Msgf("-------------url: %s, method: %s, resp time: %s", url, req.Method, time.Since(start).String())

	if 200 <= resp.StatusCode && resp.StatusCode <= 299 {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return resp, err
		}

		// restore response body with body just read
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return resp, err
}

func (r *RestClient) Get(path string, args ...any) (*Response, error) {
	queryParams := make(map[string]string)
	headers := make(map[string]string)
	if len(args) > 0 {
		queryParams = args[0].(map[string]string)
	}

	if len(args) > 1 {
		headers = args[1].(map[string]string)
	}

	request := Request{
		Method:      Get,
		BaseURL:     path,
		Headers:     headers,
		QueryParams: queryParams,
	}
	req, err := BuildRequestObject(request)
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(req)

	if err != nil {
		return nil, err
	}

	// Build Response object.
	return BuildResponse(resp)
}

// BuildRequestObject creates the HTTP request object.
func BuildRequestObject(request Request) (*http.Request, error) {
	// Add any query parameters to the URL.
	if len(request.QueryParams) != 0 {
		request.BaseURL = AddQueryParameters(request.BaseURL, request.QueryParams)
	}
	req, err := http.NewRequest(string(request.Method), request.BaseURL, bytes.NewBuffer(request.Body))
	if err != nil {
		return req, err
	}
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}
	_, exists := req.Header["Content-Type"]
	if len(request.Body) > 0 && !exists {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, err
}

// BuildResponse builds the response struct.
func BuildResponse(res *http.Response) (*Response, error) {
	body, err := io.ReadAll(res.Body)
	response := Response{
		StatusCode: res.StatusCode,
		Body:       string(body),
		Headers:    res.Header,
	}
	res.Body.Close()
	return &response, err
}

// AddQueryParameters adds query parameters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	baseURL += "?"
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + params.Encode()
}
