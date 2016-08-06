package device

import (
	"encoding/json"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
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


func MakeHTTPHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		//kithttp.ServerErrorEncoder(encodeError),
	}

	// POST		/devices/	create a device
	r.Methods("POST").Path("/devices/").Handler(httptransport.NewServer(
		ctx,
		e.PostDeviceEndpoint,
		decodePostDeviceRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodePostDeviceRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

}