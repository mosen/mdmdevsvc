package jsonapi

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"net/http"
)

type Middleware func(endpoint.Endpoint) endpoint.Endpoint

func JsonApiRequestMiddleware() Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			jsonApiRequest := request.(http.Request)
			if jsonApiRequest.Method == "POST" && jsonApiRequest.Header().Get("Content-Type") != "application/vnd.api+json" {
				jsonApiErr := &Error{
					Id: "not_acceptable",
					Status: 406,
					Title: "Not Acceptable",
					Detail: "Content-Type header must be set to 'application/vnd.api+json'",
				}

				return nil, jsonApiErr
			}

			if jsonApiRequest.Header().Get("Accept") != "application/vnd.api+json" {
				jsonApiErr := &Error{
					Id: "not_acceptable",
					Status: 406,
					Title: "Not Acceptable",
					Detail: "Accept header must be set to 'application/vnd.api+json'",
				}

				return nil, jsonApiErr
			}

			return next(ctx, request)
		}
	}
}
