package shared

import (
	"encoding/json"
	"net/http"
)

type ErrorResp struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func BadRequest(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusBadRequest, ErrorResp{Error: msg})
}

func InternalError(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusInternalServerError, ErrorResp{Error: err.Error()})
}

func NotFound(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusNotFound, ErrorResp{Error: msg})
}
