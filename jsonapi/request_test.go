package jsonapi

import (
	"testing"
	"encoding/json"
	"fmt"
)


func TestCreateRequestMarshal(t *testing.T) {
	var postRequest Request
	postRequest = Request{
		Data: data{
			Id: "1",
			Type: "devices",
			Attributes: map[string]string{
				"crud": "crud",
			},
		},
	}

	sut, err := json.Marshal(postRequest)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s\n", sut)
}
