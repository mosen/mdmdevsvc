package jsonapi

import "net/url"

type Data struct {
	Id            string          `json:"id,omitempty"`
	Type          string          `json:"type"`
	Relationships map[string]Data `json:"relationships,omitempty"`
	URL           url.URL
	Attributes    interface{} `json:"attributes,omitempty"`
}
