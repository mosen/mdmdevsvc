package jsonapi

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeData(t *testing.T) {
	reqBody := []byte(`{
			"type": "devices",
			"attributes": {
				"name": "Test Device",
				"udid": "00000000-1111-2222-3333-444455556666",
				"serial_number": "ABCDEFG0000111",
				"os_version": "9.3.3"
			}
	}`)
	var sut Data

	err := json.Unmarshal(reqBody, &sut)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", sut)
}

func TestEncodeData(t *testing.T) {
	sut := Data{
		Id: "1111-2222-3333-4444",
		Type: "devices",
	}

	jsonBytes, err := json.Marshal(&sut)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s\n", jsonBytes)

}
