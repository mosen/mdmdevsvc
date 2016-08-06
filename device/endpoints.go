package device

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/mosen/devicestore/jsonapi"
	"golang.org/x/net/context"
)

type Endpoints struct {
	PostDeviceEndpoint	endpoint.Endpoint
	GetDeviceEndpoint endpoint.Endpoint
	PutDeviceEndpoint endpoint.Endpoint
	PatchDeviceEndpoint endpoint.Endpoint
	DeleteDeviceEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostDeviceEndpoint: MakePostDeviceEndpoint(s),
	}
}

func MakePostDeviceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(jsonapi.CreateRequest)
		createDevice := req.Data.Attributes.(Device)

		objectUuid, jsonApiErr := s.PostDevice(ctx, createDevice)
		if jsonApiErr != nil {
			return jsonapi.CreateResponse{Data: nil, Errors: []jsonapi.Error{}}, nil
		}

		return jsonapi.CreateResponse{Data: &jsonapi.Data{
			Type: "devices",
			Id: objectUuid.String(),
			Attributes: req.Data.Attributes,
		}, Errors: nil}, nil
	}
}


