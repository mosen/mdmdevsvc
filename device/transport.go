package device

import (
	"encoding/json"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"github.com/mosen/devicestore/jsonapi"
)

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request jsonapi.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", request)

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/vnd.api+json")
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
		makeCreateEndpoint(svc, logger),
		decodeCreateRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/", createHandler).Methods("POST")

	return r
}
