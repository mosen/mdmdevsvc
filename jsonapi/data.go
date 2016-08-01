package jsonapi

import (
	"net/url"
	"github.com/satori/go.uuid"
)

type Data struct {
	Id            string          `json:"id,omitempty"`
	Type          string          `json:"type"`
	Attributes    interface{}     `json:"attributes,omitempty"`
}

type DataObject interface {
	Id() (*uuid.UUID, bool)
	Type() string
	Links() map[string]url.URL
	Attributes() interface{}
}
