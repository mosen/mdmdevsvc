package device

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/mosen/devicestore/jsonapi"
	"golang.org/x/net/context"
)

func makeCreateEndpoint(svc deviceService, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(jsonapi.CreateRequest)
		createDevice := req.Data.Attributes.(Device)

		objectUuid, jsonApiErr := svc.Create(&createDevice)
		if jsonApiErr != nil {
			return jsonapi.CreateResponse{Data: nil, Errors: []jsonapi.Error{*jsonApiErr}}, nil
		}

		return jsonapi.CreateResponse{Data: &jsonapi.Data{
			Type: "devices",
			Id: objectUuid.String(),
			Attributes: req.Data.Attributes,
		}, Errors: nil}, nil
	}
}

//func makeUpdateEndpoint(svc deviceService) endpoint.Endpoint {
//
//}
//
//func makeDeleteEndpoint(svc deviceService) endpoint.Endpoint {
//
//}
