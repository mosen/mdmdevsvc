package jsonapi

type Data struct {
	Type string `json:"type"`
	Attributes interface{}
	Relationships map[string]interface{}
}

type Request struct {
	Data Data `json:"data"`
}
