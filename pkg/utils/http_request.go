package utils

import "net/http"

func RequestType(r *http.Request) string {
	return r.Header.Get("Content-Type")
}

func RequestIsJson(r *http.Request) bool {
	return RequestType(r) == "application/json"
}
