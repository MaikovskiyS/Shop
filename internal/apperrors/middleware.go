package apperrors

import (
	"encoding/json"
	"net/http"
)

func ErrResponse(w http.ResponseWriter, er *AppErr) {
	resp := make(map[string]string, 1)
	w.WriteHeader(er.Code())
	resp["error"] = er.Error()

	rBytes, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(rBytes)

}
