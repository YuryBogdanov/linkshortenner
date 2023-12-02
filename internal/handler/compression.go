package handler

import (
	"compress/gzip"
	"io"
	"net/http"
)

const (
	incomingContentCompressionHeader = "Content-Encoding"
	outgoingContentCompressionHeader = "Accept-Encoding"
	encodingMethod                   = "gzip"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func withCompression(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(incomingContentCompressionHeader) != "gzip" {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set(outgoingContentCompressionHeader, encodingMethod)
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
}
