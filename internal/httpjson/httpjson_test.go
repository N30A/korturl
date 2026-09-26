package httpjson

import (
	"bytes"
	"encoding/json"
	"net/http"
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

	_, err := Decode[wantStruct](request)
	if err == nil {
		t.Error("Expected error for garbage after valid JSON, got nil")
	}
}

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()

	data := struct {
		Username string `json:"username"`
	}{Username: "test"}

	Write(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	var got struct {
		Username string `json:"username"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if got.Username != "test" {
		t.Errorf("Expected username 'test', got %s", got.Username)
	}
}

func TestWriteMarshalError(t *testing.T) {
	w := httptest.NewRecorder()

	data := struct {
		Ch chan int `json:"ch"`
	}{Ch: make(chan int)}

	Write(w, http.StatusOK, data)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
