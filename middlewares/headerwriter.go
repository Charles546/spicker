// This file is manually authored.

package middlewares

import (
	"net/http"
)

type HeaderWriter struct {
	headers         map[string]string
	w               http.ResponseWriter
	didWrite        bool
	onlyWhenSuccess bool
}

func (h *HeaderWriter) Header() http.Header {
	return h.w.Header()
}

func (h *HeaderWriter) WriteHeader(code int) {
	if !h.didWrite {
		if code == 200 || !h.onlyWhenSuccess {
			for k, v := range h.headers {
				h.w.Header().Set(k, v)
			}
		}
		h.didWrite = true
	}

	h.w.WriteHeader(code)
}

func (h *HeaderWriter) Write(b []byte) (int, error) {
	if !h.didWrite {
		for k, v := range h.headers {
			h.w.Header().Set(k, v)
		}
		h.didWrite = true
	}

	//nolint:wrapcheck
	return h.w.Write(b)
}

func (h *HeaderWriter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.w = w
		next.ServeHTTP(h, r)
		if !h.didWrite {
			for k, v := range h.headers {
				w.Header().Set(k, v)
			}
			h.didWrite = true
		}
	})
}
