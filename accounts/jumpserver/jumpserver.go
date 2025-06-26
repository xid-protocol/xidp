package jumpserver

import (
	"context"

	"github.com/colin-404/logx"
	"github.com/tidwall/gjson"
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/db/models"
	"github.com/xid-protocol/xidp/db/repositories"
)

// type jumpserver struct {
// 	accessKey string
// 	secret    string
// }

func setJumpserverInfo(userInfo gjson.Result) {
	repo := repositories.NewXidInfoRepository()
	ctx := context.Background()
	if email := userInfo.Get("email").String(); email != "" {
		xid := common.GenerateXid(email)
		logx.Infof("xid: %v", xid)
		//检查是否存在
		exists, err := repo.CheckXidInfoExists(ctx, map[string]interface{}{"xid": xid}, "/info/jumpserver")
		if err != nil {
			logx.Errorf("检查用户 %s 失败: %v", email, err)
			return
		}
		if exists {
			logx.Infof("用户 %s 已存在，跳过插入", email)
			return
		}

		xidRecord := models.XID{
			Name:    "xid-protocol",
			Xid:     xid,
			Payload: userInfo.Value(),
			Version: "0.1.0",
			Metadata: models.Metadata{
				CreatedAt:   common.GetTimestamp(),
				Operation:   "create",
				Path:        "/info/jumpserver",
				ContentType: "application/json",
			},
		}
		err = repo.Insert(ctx, &xidRecord)
		if err != nil {
			logx.Errorf("插入用户 %s 失败: %v", email, err)
		} else {
			logx.Infof("成功插入用户: %s", email)
		}
	}
}

func RunJumpServer() {

	resp := getUserInfo()
	results := gjson.Parse(resp.String())

	for _, userInfo := range results.Array() {
		setJumpserverInfo(userInfo)
	}

	// repo := repositories.NewXidInfoRepository()
	// ctx := context.Background()
	// for email, _ := range usersByEmail {
	// 	// 检查sealsuite是否已存在
	// 	info := map[string]interface{}{"email": email, "type": "user_email"}
	// 	exists, err := repo.CheckXidInfoExists(ctx, info, "/info/sealsuite")
	// 	if err != nil {
	// 		logx.Errorf("检测sealsuite失败: %v", err)
	// 		continue
	// 	}

	// 	if exists {
	// 		logx.Infof("sealsuite %s 已存在，跳过插入", email)
	// 		continue
	// 	}

	// 	now := time.Now().UTC() // 获取当前时间

	// 	xidRecord := models.XID{
	// 		Xid:         emailXid(email),
	// 		Operation:   "create",
	// 		Path:        "/info/sealsuite",
	// 		ContentType: "application/json",
	// 		Payload:     usersByEmail[email].Value(),
	// 		Version:     "0.1.0",
	// 		CreatedAt:   now,
	// 		UpdatedAt:   now,
	// 	}
	// 	logx.Infof("xid: %v", xidRecord)
	// 	// 插入MongoDB
	// 	err = repo.Insert(ctx, &xidRecord)
	// 	if err != nil {
	// 		logx.Errorf("插入用户 %s 失败: %v", email, err)
	// 	} else {
	// 		logx.Infof("成功插入用户: %s", email)
	// 	}
	// }

}
