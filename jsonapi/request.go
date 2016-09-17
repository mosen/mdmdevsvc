package jsonapi

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

type Data struct {
	Id         string             `json:"id,omitempty"`
	Type       string             `json:"type"`
	Attributes interface{}        `json:"attributes,omitempty"`
	Links      map[string]url.URL `json:"links,omitempty"`
}

type Relationship struct {
	Data Data `json:"data"`
}

// This request struct is used within POST, PUT and PATCH verbs.
type Request struct {
	Data          Data                    `json:"data"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

func DecodeJsonApiPostRequest(_ context.Context, req *http.Request) (interface{}, error) {
	fmt.Println("Decoding JSON-API Request")

	var jsonApiRequest Request
	if err := json.NewDecoder(req.Body).Decode(&jsonApiRequest); err != nil {
		return Request{}, err
	}

	fmt.Println("Decoded JSON API Request")

	return jsonApiRequest, nil
}
