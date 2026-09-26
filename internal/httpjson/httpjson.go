package httpjson

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func Write(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to encode response", "data", data, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func Decode[T any](r *http.Request) (T, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var data T
	if err := decoder.Decode(&data); err != nil {
		return data, err
	}

	if decoder.More() {
		return data, errors.New("unexpected trailing data after JSON value")
	}

	return data, nil
}
