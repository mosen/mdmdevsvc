package jsonapi

// Error source is an object describing the source of the error in the request
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// Error is a JSON-API standard error object
type Error struct {
	Id     string            `json:"id,omitempty"`
	Links  Links             `json:"links,omitempty"`
	Status string            `json:"status,omitempty"`
	Code   string            `json:"code,omitempty"`
	Title  string            `json:"title,omitempty"`
	Detail string            `json:"detail,omitempty"`
	Source ErrorSource       `json:"source,omitempty"`
	Meta   map[string]string `json:"meta,omitempty"`
}

func (e *Error) Error() string {
	return e.Detail
}
