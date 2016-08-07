package jsonapi

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

type data struct {
	Id         string             `json:"id,omitempty"`
	Type       string             `json:"type"`
	Attributes interface{}        `json:"attributes,omitempty"`
	Links      map[string]url.URL `json:"links,omitempty"`
}

type relationship struct {
	Data data `json:"data"`
}

// This request struct is used within POST, PUT and PATCH verbs.
type Request struct {
	Data          data                    `json:"data"`
	Relationships map[string]relationship `json:"relationships,omitempty"`
}

func DecodeJsonApiPostRequest(_ context.Context, req *http.Request) (interface{}, error) {
	var jsonApiRequest Request
	if err := json.NewDecoder(req.Body).Decode(&jsonApiRequest); err != nil {
		return Request{}, err
	}

	return jsonApiRequest, nil
}
