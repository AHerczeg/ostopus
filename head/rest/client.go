package rest

import (
	"net/http"
	"time"
)

func GetDefaultClient() http.Client {
	return http.Client{
		Timeout: time.Second * 10,
	}
}
