// This file is manually authored.

package app

import (
	"testing"
	"time"

	"github.com/Charles546/spicker/v2/models"
	"github.com/Charles546/spicker/v2/restapi/operations"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestFromCache(t *testing.T) {
	redisClient = &CacheStub{}
	db, mock := redismock.NewClientMock()
	redisClient.redisClient = db

	mock.ExpectGet("TEST_200").RedisNil()
	body := FromCache("TEST", 200)
	assert.Nilf(t, body, "should return nil body with non-existing key")
	mock.ClearExpect()

	serialized := `{"average":234.1,"history":[{"close":234.1,"date":"2022-04-10","high":253.78,"low":221.34,"open":223.33,"volume":3489637}],"symbol":"TEST"}`
	var av float32 = 234.1
	symbol := "TEST"
	expected := &operations.StockpricesOKBody{
		Symbol:  &symbol,
		Average: &av,
		History: []*models.Stockprice{
			&models.Stockprice{
				Date:   "2022-04-10",
				Open:   223.33,
				High:   253.78,
				Low:    221.34,
				Close:  234.1,
				Volume: 3489637,
			},
		},
	}
	mock.ExpectGet("TEST_1").SetVal(serialized)
	actual := FromCache("TEST", 1)
	assert.Equalf(t, expected, actual, "should return the stored payload with correct key")
	mock.ClearExpect()
}

func TestCache(t *testing.T) {
	redisClient = &CacheStub{}
	db, mock := redismock.NewClientMock()
	redisClient.redisClient = db

	var av float32 = 234.1
	symbol := "TST"
	payload := &operations.StockpricesOKBody{
		Symbol:  &symbol,
		Average: &av,
		History: []*models.Stockprice{
			&models.Stockprice{
				Date:   "2022-04-10",
				Open:   223.33,
				High:   253.78,
				Low:    221.34,
				Close:  234.1,
				Volume: 3489673,
			},
		},
	}
	serialized := []byte(`{"average":234.1,"history":[{"close":234.1,"date":"2022-04-10","high":253.78,"low":221.34,"open":223.33,"volume":3489637}],"symbol":"TST"}`)
	mock.CustomMatch(func(expected, actual []interface{}) error {
		assert.Equalf(t, expected[0], actual[0], "should store key TST_1")
		assert.Equalf(t, expected[1], actual[1], "should serialize the payload")
		// skip the time compare
		return nil
	}).ExpectSetEX("TST_1", []byte(serialized), time.Second).SetVal("1")
	Cache("TST", 1, payload)
	mock.ClearExpect()
}
