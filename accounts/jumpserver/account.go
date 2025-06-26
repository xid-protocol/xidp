package jumpserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/colin-404/logx"
	"github.com/go-fed/httpsig"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/xid-protocol/xidp/common"
)

func JumpServerAccountMonitor() *resty.Response {
	url := viper.GetString("jumpserver.endpoint") + "/api/v1/users/users/"

	now := time.Now().UTC()
	gmtDate := now.Format("Mon, 02 Jan 2006 15:04:05 GMT")

	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.HMAC_SHA256},
		httpsig.DigestSha256,
		[]string{"(request-target)", "accept", "date"},
		httpsig.Signature,
		0,
	)
	if err != nil {
		return nil
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	req.Header.Set("Date", gmtDate)

	// 创建HMAC密钥
	secret := []byte(viper.GetString("jumpserver.access_key_secret"))
	mac := hmac.New(sha256.New, secret)

	// 修正SignRequest调用
	err = signer.SignRequest(
		mac, // crypto.PrivateKey (HMAC)
		viper.GetString("jumpserver.access_key_id"), // keyId
		req, // *http.Request
		nil, // body ([]byte)
	)
	if err != nil {
		return nil
	}

	// 发送请求
	resp := common.DoHttp("GET", url, nil, map[string]string{
		"Accept":        req.Header.Get("Accept"),
		"X-JMS-ORG":     req.Header.Get("X-JMS-ORG"),
		"Date":          req.Header.Get("Date"),
		"Authorization": req.Header.Get("Authorization"),
	})
	logx.Infof("resp: %v", resp.String())
	logx.Infof("resp: %v", resp.StatusCode())
	return resp
}
