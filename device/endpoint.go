package device

import (
	"golang.org/x/net/context"
	"github.com/go-kit/kit/endpoint"
)

type createRequest struct {
	data Device `json:"data"`
}

func MakeCreateEndpoint(svc deviceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

	}
}