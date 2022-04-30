// This file is manually authored.

package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpires(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost/stockprices", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>test</body></html>")
	})
	Expires(handler).ServeHTTP(w, r)

	resp := w.Result()
	assert.NotEmptyf(t, resp.Header.Get("Expires"), "should set expires header on 200 result")

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "http://localhost/stockprices", nil)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})
	Expires(handler).ServeHTTP(w, r)

	resp = w.Result()
	assert.Emptyf(t, resp.Header.Get("Expires"), "should NOT set expires header on non-200 result")
}
