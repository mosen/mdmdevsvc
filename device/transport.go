package device

import (
	"encoding/json"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"fmt"
)

type Data struct {
	Id         string             `json:"id,omitempty"`
	Type       string             `json:"type"`
	Attributes *Device            `json:"attributes,omitempty"`
	Links      map[string]url.URL `json:"links,omitempty"`
}

type Relationship struct {
	Data Data `json:"data"`
}

// This request struct is used within POST, PUT and PATCH verbs.
type Request struct {
	Data          Data                    `json:"data"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	return json.NewEncoder(w).Encode(response)
}

func decodePostDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req Request
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodePostDeviceResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if device, ok := response.(*Device); ok {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("Location", fmt.Sprintf("/devices/%s", device.UUID))

		var jsonApiResponse Request
		jsonApiResponse = Request{
			Data: Data{
				Id: device.UUID.String(),
				Type: "devices",
				Attributes: device,
			},
		}
		fmt.Printf("%v\n", jsonApiResponse)

		return json.NewEncoder(w).Encode(jsonApiResponse)
	} else {
		fmt.Println("response was not device")
	}
	return nil
}

func MakeHTTPHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST		/devices/	create a device
	r.Methods("POST").Path("/devices").Handler(httptransport.NewServer(
		ctx,
		e.PostDeviceEndpoint,
		decodePostDeviceRequest,
		encodePostDeviceResponse,
		options...,
	))

	return r
}


func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	//case ErrAlreadyExists, ErrInconsistentIDs:
	//	return http.StatusBadRequest
	default:
		if e, ok := err.(httptransport.Error); ok {
			switch e.Domain {
			case httptransport.DomainDecode:
				return http.StatusBadRequest
			case httptransport.DomainDo:
				return http.StatusServiceUnavailable
			default:
				return http.StatusInternalServerError
			}
		}
		return http.StatusInternalServerError
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}