package jsonapi

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestCreateRequestUnmarshal(t *testing.T) {
	var jsonBytes []byte = []byte(`{"data":{"id":"12345","type":"devices"}}`)
	var sut CreateRequest

	if err := json.Unmarshal(jsonBytes, &sut); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", sut)
}

func TestCreateRequestMarshal(t *testing.T) {
	var cr CreateRequest = CreateRequest{
		Data{
			Id: "12345",
			Type: "devices",
		},
	}

	sut, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s\n", sut)
}
