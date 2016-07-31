package device

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/mosen/devicestore/jsonapi"
	"golang.org/x/net/context"
)

type createRequest struct {
	jsonapi.Data
}

type createResponse struct {
	device *Device `json:"data,omitempty"`
	err    *jsonapi.Error
}

func makeCreateEndpoint(svc deviceService, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(createRequest)
		logger.Log("Creating new device")
		_, jsonApiErr := svc.Create(&Device{})
		if jsonApiErr != nil {
			//return createResponse{device: nil, err: jsonApiErr}, nil
		}
		//
		//return createResponse{device: &req.data.Attributes, err: nil}, nil
		return nil, nil
	}
}

//func makeUpdateEndpoint(svc deviceService) endpoint.Endpoint {
//
//}
//
//func makeDeleteEndpoint(svc deviceService) endpoint.Endpoint {
//
//}
