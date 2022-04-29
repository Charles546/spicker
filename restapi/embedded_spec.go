// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Stock price grabber",
    "version": "0.1.0"
  },
  "paths": {
    "/stockprices": {
      "get": {
        "produces": [
          "application/json"
        ],
        "operationId": "stockprices",
        "responses": {
          "200": {
            "description": "return certain history prices and the average closing prices for a certain stock.",
            "schema": {
              "description": "the stock prices",
              "type": "object",
              "required": [
                "history",
                "average",
                "symbol"
              ],
              "properties": {
                "average": {
                  "description": "average closing price of the past N days",
                  "type": "number",
                  "format": "float",
                  "example": 15.56
                },
                "history": {
                  "description": "the history prices",
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/stockprice"
                  }
                },
                "symbol": {
                  "description": "the stock symbol",
                  "type": "string",
                  "example": "QQQ"
                }
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "stockprice": {
      "type": "object",
      "properties": {
        "close": {
          "description": "closing price",
          "type": "number",
          "format": "float"
        },
        "date": {
          "description": "the date of the prices",
          "type": "string"
        },
        "high": {
          "description": "highest price of the day",
          "type": "number",
          "format": "float"
        },
        "low": {
          "description": "lowest price of the day",
          "type": "number",
          "format": "float"
        },
        "open": {
          "description": "opening price",
          "type": "number",
          "format": "float"
        },
        "volume": {
          "description": "traded volume",
          "type": "integer",
          "format": "int64"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Stock price grabber",
    "version": "0.1.0"
  },
  "paths": {
    "/stockprices": {
      "get": {
        "produces": [
          "application/json"
        ],
        "operationId": "stockprices",
        "responses": {
          "200": {
            "description": "return certain history prices and the average closing prices for a certain stock.",
            "schema": {
              "description": "the stock prices",
              "type": "object",
              "required": [
                "history",
                "average",
                "symbol"
              ],
              "properties": {
                "average": {
                  "description": "average closing price of the past N days",
                  "type": "number",
                  "format": "float",
                  "example": 15.56
                },
                "history": {
                  "description": "the history prices",
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/stockprice"
                  }
                },
                "symbol": {
                  "description": "the stock symbol",
                  "type": "string",
                  "example": "QQQ"
                }
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "stockprice": {
      "type": "object",
      "properties": {
        "close": {
          "description": "closing price",
          "type": "number",
          "format": "float"
        },
        "date": {
          "description": "the date of the prices",
          "type": "string"
        },
        "high": {
          "description": "highest price of the day",
          "type": "number",
          "format": "float"
        },
        "low": {
          "description": "lowest price of the day",
          "type": "number",
          "format": "float"
        },
        "open": {
          "description": "opening price",
          "type": "number",
          "format": "float"
        },
        "volume": {
          "description": "traded volume",
          "type": "integer",
          "format": "int64"
        }
      }
    }
  }
}`))
}