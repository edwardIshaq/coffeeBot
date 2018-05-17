package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware middleWare will try to gzip the responses
type GzipMiddleware struct {
	Next http.Handler
}

func (mw *GzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mw.Next == nil {
		mw.Next = http.DefaultServeMux
	}

	encodings := r.Header.Get("Accept-Encoding")
	if !strings.Contains(encodings, "gzip") {
		mw.Next.ServeHTTP(w, r)
		return
	}
	w.Header().Add("Content-Encoding", "gzip")
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()
	grw := gzipResponseWriter{
		ResponseWriter: w,
		Writer:         gzipWriter,
	}
	mw.Next.ServeHTTP(grw, r)
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (grw gzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}
