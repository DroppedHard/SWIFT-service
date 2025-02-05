package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/utils"
)

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil || r.Body == http.NoBody {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteMessage(w http.ResponseWriter, status int, mess string) {
	WriteJson(w, status, map[string]string{utils.ResponseMessageField: mess})
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{utils.ResponseMessageField: err.Error()})
}
