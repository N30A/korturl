package httpjson

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "test",
		Password: "test",
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	request := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	type wantStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	got, err := Decode[wantStruct](request)
	if err != nil {
		t.Errorf("Wanted nil, got %v", err)
	}

	want := wantStruct{
		Username: "test",
		Password: "test",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestDecodeDisallowedUnknownFields(t *testing.T) {
	data := struct {
		UnknownField string `json:"unknown_field"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}{
		UnknownField: "unknown_field",
		Username:     "test",
		Password:     "test",
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	request := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	type wantStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	_, err = Decode[wantStruct](request)
	if err == nil {
		t.Error("Expected error for unknown field, got nil")
	}
}

func TestDecodeUnknownData(t *testing.T) {
	request := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("some data")))

	type wantStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	_, err := Decode[wantStruct](request)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDecodeEmptyBodyReaderNil(t *testing.T) {
	request := httptest.NewRequest("POST", "/", bytes.NewReader(nil))

	type wantStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	_, err := Decode[wantStruct](request)
	if err == nil {
		t.Error("Expected error for empty body, got nil")
	}
}

func TestDecodeWithGarbage(t *testing.T) {
	data := `
		{
			"username": "test",
			"password": "test"
		} garbage
	`
	request := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(data)))

	type wantStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ehm, err := Decode[wantStruct](request)
	t.Log(err)
	t.Log(ehm)
	if err == nil {
		t.Error("Expected error for garbage after valid JSON, got nil")
	}
}
