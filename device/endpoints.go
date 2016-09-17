package device

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Endpoints struct {
	PostDeviceEndpoint   endpoint.Endpoint
	GetDeviceEndpoint    endpoint.Endpoint
	PutDeviceEndpoint    endpoint.Endpoint
	PatchDeviceEndpoint  endpoint.Endpoint
	DeleteDeviceEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostDeviceEndpoint: MakePostDeviceEndpoint(s),
	}
}

func MakePostDeviceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)

		var d *Device = req.Data.Attributes
		if err := s.PostDevice(ctx, d); err != nil {
			return nil, err
		}

		fmt.Printf("returning %T", d)

		return d, nil
	}
}

type postDeviceRequest struct {
}
