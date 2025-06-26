// Golang 示例
package jumpserver

import (
	"log"
	"net/http"
	"time"

	"github.com/colin-404/logx"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/xid-protocol/xidp/common"
	"gopkg.in/twindagger/httpsig.v1"
)

type SigAuth struct {
	KeyID    string
	SecretID string
}

func (auth *SigAuth) Sign(r *http.Request) error {
	headers := []string{"(request-target)", "date"}
	signer, err := httpsig.NewRequestSigner(auth.KeyID, auth.SecretID, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(r, headers, nil)
}

func userInfo(jmsurl string, auth *SigAuth) *resty.Response {
	url := jmsurl + "/api/v1/users/users/"
	logx.Infof("url: %v", url)
	// logx.Infof("auth: %v", auth.KeyID)
	// logx.Infof("auth: %v", auth.SecretID)
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	// client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	if err != nil {
		log.Fatal(err)
	}
	if err := auth.Sign(req); err != nil {
		log.Fatal(err)
	}

	resp := common.DoHttp("GET", url, nil, map[string]string{
		"Accept":        req.Header.Get("Accept"),
		"X-JMS-ORG":     req.Header.Get("X-JMS-ORG"),
		"Date":          req.Header.Get("Date"),
		"Authorization": req.Header.Get("Authorization"),
	})

	return resp
}

func getUserInfo() *resty.Response {
	auth := SigAuth{
		KeyID:    viper.GetString("jumpserver.access_key_id"),
		SecretID: viper.GetString("jumpserver.access_key_secret"),
	}
	logx.Infof("auth: %v", auth)
	return userInfo(viper.GetString("jumpserver.endpoint"), &auth)
}
