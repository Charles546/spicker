// This file is manually authored.

package restapi

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "golang.org/x/exp/maps"
	"net/http"
    "os"
    "sort"
    "strconv"

	"github.com/go-openapi/runtime/middleware"

	"github.com/Charles546/spicker/v2/models"
	"github.com/Charles546/spicker/v2/restapi/operations"
)

const APIPattern string = "https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s&outputsize=%s"

func getStockprices(params operations.StockpricesParams) middleware.Responder {
    symbol := os.Getenv("SYMBOL")
    ndaysStr := os.Getenv("NDAYS")
    apikey := os.Getenv("ALPHAVANTAGE_APIKEY")
    if len(symbol)==0 {
        return middleware.Error(http.StatusInternalServerError, "configuration missing for operations.Stockprices: SYMBOL")
    }
    if len(ndaysStr)==0 {
        return middleware.Error(http.StatusInternalServerError, "configuration missing for operations.Stockprices: NDAYS")
    }
    if len(apikey)==0 {
        return middleware.Error(http.StatusInternalServerError, "configuration missing for operations.Stockprices: ALPHAVANTAGE_APKEY")
    }
    ndays, err := strconv.Atoi(ndaysStr)
    if err!=nil {
        return middleware.Error(http.StatusInternalServerError, "invalid configuration for operations.Stockprices: NDAYS")
    }
    outputSize := "compact"
    if ndays>100 {
        outputSize = "full"
    }

    url := fmt.Sprintf(APIPattern, apikey, symbol, outputSize);
    resp, err := http.Get(url)

    if err!=nil {
        log.Printf("operations.Stockprices api response error: %w\n", err);
        return middleware.Error(http.StatusInternalServerError, "operations.Stockprices failed")
    }

    var j map[string]interface{}
    body, err := ioutil.ReadAll(resp.Body)
    if err==nil {
        json.Unmarshal(body, &j)
    }
    if err!=nil {
        log.Printf("operations.Stockprices api parse error: %w\n", err);
        return middleware.Error(http.StatusInternalServerError, "operations.Stockprices failed")
    }

    hist := j["Time Series (Daily)"].(map[string]interface{})
    dates := maps.Keys(hist)
    sort.Slice(dates, func(i, j int) bool { return dates[i] > dates[j] })

    ret := operations.StockpricesOKBody{}
    ret.Symbol = &symbol
    ret.History = make([]*models.Stockprice, ndays)

    var sum float32
    count := 0
    for i := range dates {
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

        ret.History[count] = &day;

        sum += day.Close
        count++
        if(count==ndays) {
            average := sum / float32(ndays)
            ret.Average = &average
            break
        }
    }

    return operations.NewStockpricesOK().WithPayload(&ret)
}
