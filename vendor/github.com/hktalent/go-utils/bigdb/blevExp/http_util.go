package blevExp

//go:generate go-bindata-assetfs -pkg=main ../static/...
//go:generate go fmt .

import (
	util "github.com/hktalent/go-utils"
	"io"
	"log"
	"net/http"
	"strings"
)

type MyFileHandler struct {
	H http.Handler
}

func (mfh MyFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if *StaticEtag != "" {
		w.Header().Set("Etag", *StaticEtag)
	}
	mfh.H.ServeHTTP(w, r)
}

func RewriteURL(to string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = to
		h.ServeHTTP(w, r)
	})
}

func muxVariableLookup(req *http.Request, name string) string {
	return ""
}

func DocIDLookup(req *http.Request) string {
	a := strings.Split(req.URL.Path, "/")
	if 2 < len(a) {
		return a[3]
	}
	if "" != req.FormValue("docID") {
		return req.FormValue("docID")
	}
	return muxVariableLookup(req, "docID")
}

func IndexNameLookup(req *http.Request) string {
	a := strings.Split(req.URL.Path, "/")
	if 1 < len(a) {
		return a[2]
	}
	return muxVariableLookup(req, "indexName")
}

func showError(w http.ResponseWriter, r *http.Request,
	msg string, code int) {
	log.Printf("Reporting error %v/%v", code, msg)
	http.Error(w, msg, code)
}

func mustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := util.Json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}
