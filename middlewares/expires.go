// This file is manually authored.

package middlewares

import (
	"net/http"
	"time"
)

func Expires(next http.Handler) http.Handler {
	newYork, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(newYork)
	available := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, newYork)
	if available.Before(now) {
		available.AddDate(0, 0, 1)
	}

	expires := available.In(time.UTC).Format("Mon, 2 Jan 2006 15:04:05 MST")

	hw := &HeaderWriter{
		headers: map[string]string{
			"Expires": expires,
		},
		onlyWhenSuccess: true,
	}

	return hw.Handler(next)
}
