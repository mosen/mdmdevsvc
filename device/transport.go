package device

import (
	"encoding/json"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
)

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// ServiceHandler returns an HTTP Handler for the devices service
func ServiceHandler(ctx context.Context, svc deviceService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	r := mux.NewRouter()

	createHandler := kithttp.NewServer(
		ctx,
		makeCreateEndpoint(svc),
		decodeCreateRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/", createHandler).Methods("POST")

	return r
}
