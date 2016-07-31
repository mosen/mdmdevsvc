package device

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/mosen/devicestore/jsonapi"
	"golang.org/x/net/context"
)

type createRequest struct {
	device Device `json:"data"`
}

type createResponse struct {
	device *Device `json:"data,omitempty"`
	err    *jsonapi.Error
}

func makeCreateEndpoint(svc deviceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createRequest)
		_, jsonApiErr := svc.Create(&req.device)
		if jsonApiErr != nil {
			return createResponse{device: nil, err: jsonApiErr}, nil
		}

		return createResponse{device: &req.device, err: nil}, nil
	}
}

//func makeUpdateEndpoint(svc deviceService) endpoint.Endpoint {
//
//}
//
//func makeDeleteEndpoint(svc deviceService) endpoint.Endpoint {
//
//}
