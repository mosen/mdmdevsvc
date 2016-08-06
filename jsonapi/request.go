package jsonapi

import (

)

type CreateRequest struct {
	Data Data `json:"data"`
}

//func DecodeJSONApiCreateRequest(_ context.Context, req *http.Request) (request interface{}, error) {
//
//}