package TempoZaiko0x0

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 今回はすべてjsonで返却
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}
