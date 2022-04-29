// This file is manually authored.

package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/Charles546/spicker/v2/models"
	"github.com/Charles546/spicker/v2/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/exp/maps"
)

// StockpricesGrabberConfig is a data structure storing the configs.
type StockpricesGrabberConfig struct {
	symbol string
	ndays  int
	url    string
	err    middleware.Responder
}

// APIConfig is the loaded configs.
var APIConfig *StockpricesGrabberConfig

// APIPattern is a string pattern for creating the url.
const APIPattern string = "https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s&outputsize=%s"

func getConfig() *StockpricesGrabberConfig {
	if APIConfig != nil {
		return APIConfig
	}

	APIConfig = &StockpricesGrabberConfig{}
	APIConfig.symbol = os.Getenv("SYMBOL")
	ndaysStr := os.Getenv("NDAYS")
	apikey := os.Getenv("ALPHAVANTAGE_APIKEY")

	if len(APIConfig.symbol) == 0 {
		APIConfig.err = middleware.Error(http.StatusInternalServerError, "configuration missing for operations.Stockprices: SYMBOL")

		return APIConfig
	}
	if len(ndaysStr) == 0 {
		APIConfig.err = middleware.Error(http.StatusInternalServerError, "configuration missing for operations.Stockprices: NDAYS")

		return APIConfig
	}
	if len(apikey) == 0 {
		APIConfig.err = middleware.Error(http.StatusInternalServerError, "configuration missing for operations.Stockprices: ALPHAVANTAGE_APKEY")

		return APIConfig
	}

	var err error
	APIConfig.ndays, err = strconv.Atoi(ndaysStr)
	if err != nil {
		APIConfig.err = middleware.Error(http.StatusInternalServerError, "invalid configuration for operations.Stockprices: NDAYS")

		return APIConfig
	}
	outputSize := "compact"
	if APIConfig.ndays > 100 {
		outputSize = "full"
	}

	APIConfig.url = fmt.Sprintf(APIPattern, apikey, APIConfig.symbol, outputSize)

	return APIConfig
}

func getStockprices(params operations.StockpricesParams) middleware.Responder {
	c := getConfig()
	if c.err != nil {
		return c.err
	}

	resp, err := http.Get(c.url)
	if err != nil {
		log.Printf("operations.Stockprices api response error: %+v\n", err)

		return middleware.Error(http.StatusInternalServerError, "operations.Stockprices failed")
	}

	var j map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &j)
	}
	if err != nil {
		log.Printf("operations.Stockprices api parse error: %+v\n", err)

		return middleware.Error(http.StatusInternalServerError, "operations.Stockprices failed")
	}

	hist := j["Time Series (Daily)"].(map[string]interface{})
	dates := maps.Keys(hist)
	sort.Slice(dates, func(i, j int) bool { return dates[i] > dates[j] })

	ret := operations.StockpricesOKBody{}
	ret.Symbol = &c.symbol
	ret.History = make([]*models.Stockprice, c.ndays)

	var sum float32
	for i := 0; i < c.ndays; i++ {
		prices := hist[dates[i]].(map[string]interface{})
		open, _ := strconv.ParseFloat(prices["1. open"].(string), 32)
		high, _ := strconv.ParseFloat(prices["2. high"].(string), 32)
		low, _ := strconv.ParseFloat(prices["3. low"].(string), 32)
		cl, _ := strconv.ParseFloat(prices["4. close"].(string), 32)

		day := models.Stockprice{}
		day.Date = dates[i]
		day.Volume, _ = strconv.ParseInt(prices["5. volume"].(string), 10, 64)
		day.Open = float32(open)
		day.High = float32(high)
		day.Low = float32(low)
		day.Close = float32(cl)

		ret.History[i] = &day

		sum += day.Close
	}
	average := sum / float32(c.ndays)
	ret.Average = &average

	return operations.NewStockpricesOK().WithPayload(&ret)
}
