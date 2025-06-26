package sealsuite

import (
	"context"

	"github.com/colin-404/logx"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/db/models"
	"github.com/xid-protocol/xidp/db/repositories"
)

type sealsuite struct {
	endpoint string
	token    string
}

func setSealsuiteInfo(usersByEmail map[string]gjson.Result) {
	repo := repositories.NewXidInfoRepository()
	ctx := context.Background()
	for email := range usersByEmail {
		// 检查sealsuite是否已存在
		info := map[string]interface{}{"email": email, "type": "user_email"}
		exists, err := repo.CheckXidInfoExists(ctx, info, "/info/sealsuite")
		if err != nil {
			logx.Errorf("检测sealsuite失败: %v", err)
			continue
		}

		if exists {
			logx.Infof("sealsuite %s 已存在，跳过插入", email)
			continue
		}

		xidRecord := models.XID{
			Name:    "xid-protocol",
			Xid:     common.GenerateXid(email),
			Version: "0.1.0",
			Payload: usersByEmail[email].Value(),
			Metadata: models.Metadata{
				CreatedAt:   common.GetTimestamp(),
				Operation:   "create",
				Path:        "/info/sealsuite",
				ContentType: "application/json",
			},
		}
		logx.Infof("xid: %v", xidRecord)
		// 插入MongoDB
		err = repo.Insert(ctx, &xidRecord)
		if err != nil {
			logx.Errorf("插入用户 %s 失败: %v", email, err)
		} else {
			logx.Infof("成功插入用户: %s", email)
		}
	}
}

func RunSealsuite() {
	users := SealsuiteAccount()
	logx.Infof("users: %v", users)
	//获取所有所有邮箱并去重，struct{}{}

	usersByEmail := make(map[string]gjson.Result)
	for _, user := range *users {
		if user.Get("email").Exists() {
			email := user.Get("email").String()
			if email != "" {
				// 如果email已存在，会自动覆盖（去重）
				usersByEmail[email] = user

			}
		}
	}
	// emailXidInit(usersByEmail)
	setSealsuiteInfo(usersByEmail)
}

func SealsuiteAccount() *[]gjson.Result {
	endpoint := viper.GetString("sealsuite.endpoint")
	accessKeyId := viper.GetString("sealsuite.access_key_id")
	accessKeySecret := viper.GetString("sealsuite.access_key_secret")

	ss := &sealsuite{
		endpoint: endpoint,
		token:    getToken(endpoint, accessKeyId, accessKeySecret),
	}
	logx.Infof("ss: %v", ss)
	departments := ss.getDeparment()
	return ss.getAllUsersForDepartments(departments)

}

func getToken(endpoint string, accessKeyId string, accessKeySecret string) string {
	url := endpoint + "/api/open/v1/token"
	body := map[string]string{
		"access_key_id":     accessKeyId,
		"access_key_secret": accessKeySecret,
	}
	resp := common.DoHttp("POST", url, body, nil)
	token := gjson.Get(resp.String(), "data.access_token").String()

	return token

}
