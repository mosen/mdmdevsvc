package device

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type createRequest struct {
	data Device `json:"data"`
}

func makeCreateEndpoint(svc deviceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

	}
}

func makeUpdateEndpoint(svc deviceService) endpoint.Endpoint {

}

func makeDeleteEndpoint(svc deviceService) endpoint.Endpoint {

}
