package shared

import "net/http"

func WriteResponse(w http.ResponseWriter, code int, response []byte) {
	w.WriteHeader(code)
	if len(response) > 0 {
		w.Write(response)
	}
}
