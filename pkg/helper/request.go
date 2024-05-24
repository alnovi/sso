package helper

import "net/http"

func RequestIsJson(r *http.Request) bool {
	isAppJson := r.Header.Get("Content-Type") == "application/json"
	isXmlHttp := r.Header.Get("X-Requested-With") == "XMLHttpRequest"
	return isAppJson || isXmlHttp
}
