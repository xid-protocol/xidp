package common

import (
	"errors"
	"time"

	"github.com/colin-404/logx"
	"github.com/go-resty/resty/v2"
)

var httpClient = resty.New().SetTimeout(10 * time.Second)

func DoHttp(method, url string, body interface{}, headers map[string]string) (*resty.Response, error) {
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
		logx.Errorf("Unsupported HTTP method: %s", method)
		return nil, errors.New("unsupported http method")
	}

	if err != nil {
		logx.Errorf("DoHttp: %v", err)
		return nil, err
	}

	return resp, nil
}
