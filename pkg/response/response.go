package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Page      int    `json:"page,omitempty"`
	Total     int    `json:"total,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, statusCode int, data any) {
	JSON(w, statusCode, APIResponse{Success: true, Data: data})
}

func Error(w http.ResponseWriter, statusCode int, errMsg string) {
	JSON(w, statusCode, APIResponse{Success: false, Error: errMsg})
}
