package common

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

var httpClient = resty.New().SetTimeout(5 * time.Second)

func DoHttp(method, url string, body interface{}, headers map[string]string) *resty.Response {
	req := httpClient.R()

	//设置请求的token
	// if token != "" {
	// 	req.SetAuthToken(token)
	// }

	//设置header
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(url)
	case "POST":
		if body != nil {
			req.SetBody(body)
		}
		resp, err = req.Post(url)
	// Add other methods as needed
	default:
		log.Printf("Unsupported HTTP method: %s", method)
		return nil
	}

	if err != nil {
		log.Println(err)
		return nil
	}

	return resp
}
