package blevExp

import "net/http"

type ServeHTTPFace interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
