package whitelist

import (
	"github.com/spf13/viper"
	"github.com/xid-protocol/xidp/common"
	"github.com/xid-protocol/xidp/protocols"
)

func NewWhitelist(xid string, operation string, payload interface{}) *protocols.XID {
	xiddata := protocols.XID{
		Name:    viper.GetString("xid.name"),
		Xid:     xid,
		Version: viper.GetString("xid.version"),
		Metadata: &protocols.Metadata{
			CreatedAt:   common.GetTimestamp(),
			CardId:      common.GenerateCardId(),
			Operation:   operation,
			Path:        "/protocol/whitelist",
			ContentType: "application/json",
		},
		Payload: payload,
	}
	return &xiddata
}
