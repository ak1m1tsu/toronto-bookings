package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func JSON(writer http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	writer.Write(buf.Bytes())
}
