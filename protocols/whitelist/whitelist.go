package whitelist

import (
	"context"
	"errors"

	"github.com/colin-404/logx"
	"github.com/xid-protocol/xidp/db"
	"github.com/xid-protocol/xidp/protocols"
)

type AWSOpenPort struct {
	InstanceID string `json:"instanceId"`
	Cidr       string `json:"cidr"`
	FromPort   int    `json:"fromPort"`
	ToPort     int    `json:"toPort"`
	Protocol   string `json:"protocol"`
}

type Whitelist struct {
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	Sha256Value string      `json:"sha256Value"`
}

func NewWhitelist(xid string, whitelistType string, payload Whitelist) (*protocols.XID, error) {

	whitelistRepository := db.NewXidInfoRepository()
	xidInfoRepository := db.NewXidInfoRepository()
	xidInfo, err := xidInfoRepository.FindByName(context.Background(), xid, "/protocols/whitelist")
	if err != nil {
		logx.Errorf("NewWhitelist: %v", err)
		return nil, err
	}

	metadata := protocols.NewMetadata("create", "/protocols/whitelist", "application/json")

	NewXID := protocols.NewXID(xidInfo.Info, &metadata, payload)

	//检测数据库中是否存在相同的whitelist
	xidInfos, err := xidInfoRepository.FindAllByXid(context.Background(), xid)
	if err != nil {
		logx.Errorf("NewWhitelist: %v", err)
		return nil, err
	}
	if len(xidInfos) > 0 {
		//检测是否有重复的
		for _, xidInfo := range xidInfos {
			//对比Sha256Value
			if xidInfo.Payload.(Whitelist).Sha256Value == payload.Sha256Value {
				logx.Errorf("Whitelist already exists: %v", err)
				return xidInfo, errors.New("Whitelist already exists")
			}
		}
	}

	whitelistRepository.Insert(context.Background(), NewXID)

	return NewXID, nil
}
