package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
	Reverse(string) (string, error)
	Truncate(string, int) (string, error)
}

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func (stringService) Reverse(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	r := []rune(s)
	for i, j := 0, len(r) - 1; i < len(r) / 2; i, j = i + 1, j - 1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

func (stringService) Truncate(s string, l int) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	//todo - l error checking

	return s[:l], nil
}

func main() {
	ctx := context.Background()
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		ctx,
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	reverseHandler := httptransport.NewServer(
		ctx,
		makeReverseEndpoint(svc),
		decodeReverseRequest,
		encodeResponse,
	)

	truncateHandler := httptransport.NewServer(
		ctx,
		makeTruncateEndpoint(svc),
		decodeTruncateRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/reverse", reverseHandler)
	http.Handle("/truncate", truncateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(stringRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(stringRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

func makeReverseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(stringRequest)
		v, err := svc.Reverse(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeTruncateEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(truncateRequest)
		v, err := svc.Truncate(req.S, req.L)
		if err != nil {
			return truncateResponse{v, err.Error()}, nil
		}
		return truncateResponse{v, ""}, nil
	}
}

func decodeUppercaseRequest(r *http.Request, clazz interface{}) (interface{}, error) {
	var request stringRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(r *http.Request) (interface{}, error) {
	var request stringRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeReverseRequest(r *http.Request) (interface{}, error) {
	var request stringRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeTruncateRequest(r *http.Request) (interface{}, error) {
	var request truncateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type stringRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type countResponse struct {
	V int `json:"v"`
}

type reverseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type truncateRequest struct {
	S string `json:"s"`
	L int `json:"l"`
}

type truncateResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")