// This file is manually authored.

package restapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Charles546/spicker/v2/restapi/operations"
	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
)

func TestGetStockprices(t *testing.T) {
	os.Setenv("REDIS_CONNECTION", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"Time Series (Daily)":{"2022-04-29":{"1. open":"12.12","2. high":"15.99","3. low":"12.09","4. close":"13.46","5. volume":"235233"}}}`)
	}))
	defer server.Close()
	APIConfig = &StockpricesGrabberConfig{
		symbol: "MOCK",
		ndays:  1,
		url:    server.URL,
	}
	defer func() { APIConfig = nil }()

	expected := `{"average":13.46,"history":[{"close":13.46,"date":"2022-04-29","high":15.99,"low":12.09,"open":12.12,"volume":235233}],"symbol":"MOCK"}` + "\n"

	w := httptest.NewRecorder()
	os.Setenv("REDIS_CONNECTION", "")
	getStockprices(operations.StockpricesParams{}).WriteResponse(w, runtime.JSONProducer())
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equalf(t, expected, string(body), "should return the processed prices")
}
