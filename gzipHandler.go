package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzHandler struct {
	path string
}

func (g *GzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := http.FileServer(http.Dir(g.path))
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.ServeHTTP(w, r)
		return
	}
	gzpw := gzip.NewWriter(w)
	defer gzpw.Close()
	w.Header().Set("Content-Encoding", "gzip")
	gzipWriter := &gzipResponseWriter{
		ResponseWriter: w,
		Writer:         gzpw,
	}
	h.ServeHTTP(gzipWriter, r)
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (g gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}
