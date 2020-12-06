package renderer

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, response interface{}) error {
	body, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(body)
	return err
}
