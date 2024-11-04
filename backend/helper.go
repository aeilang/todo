package main

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, statusCode int, err error) error {
	return SendJSON(w, statusCode, map[string]string{
		"msg": err.Error(),
	})
}
