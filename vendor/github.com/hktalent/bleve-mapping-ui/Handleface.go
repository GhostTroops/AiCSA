package mapping

import "net/http"

type HandleFuncFace interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) MethodsFace
}

type MethodsFace interface {
	Methods(methods ...string)
}
