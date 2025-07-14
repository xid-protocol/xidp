package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/xid-protocol/xidp/db/repositories"
	"github.com/xid-protocol/xidp/protocols"
)

func GetXid(username string, source string) (*protocols.XID, error) {
	ctx := context.Background()
	path := "/info"
	if source != "" {
		xidInfoRepository := repositories.NewXidInfoRepository()
		path = fmt.Sprintf("/info/%s", source)
		xid, err := xidInfoRepository.FindByName(ctx, username, path)
		if err != nil {
			return nil, err
		}
		return xid, nil
	}

	// if path == "/info" {
	// 	xidRepository := repository.NewXIDRepository()
	// 	xid, err := xidRepository.FindByName(context.Background(), username)
	// 	if err != nil {
	// 		return ""
	// 	}
	// 	return xid.Xid
	// } else {
	// 	xidInfoRepository := repository.NewXidInfoRepository()
	// 	xid, err := xidInfoRepository.FindByXidInfo(context.Background(), username, path)
	// 	if err != nil {
	// 		logx.Errorf("Failed to find xid by name: %v", err)
	// 		return err.Error()
	// 	}
	// 	return xid.Xid
	// }

	// orConditions := []bson.M{
	// 	{"payload.username": bson.M{"$regex": username, "$options": "i"}},
	// 	{"payload.email": bson.M{"$regex": username, "$options": "i"}},
	// 	{"payload.name": bson.M{"$regex": username, "$options": "i"}},
	// }

	return nil, errors.New("xid not found")
}
