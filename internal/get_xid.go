package internal

import (
	"errors"
	"fmt"

	"github.com/xid-protocol/xidp/protocols"
)

// 转换XIDType
func ConvertXIDInfo(info map[string]interface{}) (protocols.Info, error) {
	var XIDInfo protocols.Info

	//必须要有id
	var ok bool
	if XIDInfo.ID, ok = info["id"].(string); !ok {
		return XIDInfo, fmt.Errorf("info.id is required")
	}

	//必须要有type
	if XIDInfo.Type, ok = info["type"].(string); !ok {
		return XIDInfo, fmt.Errorf("info.type is required")
	}

	return XIDInfo, nil
}

// MapToMetadata converts a generic map (e.g., JSON body) to a Metadata struct.
// It validates that required string fields (operation, path, contentType) are present.
// Returns an error if any required field is missing or not a string.
func MapToMetadata(m map[string]interface{}) (protocols.Metadata, error) {
	var md protocols.Metadata
	var ok bool
	if m == nil {
		return md, fmt.Errorf("metadata is required")
	}

	// 必须有 path
	if md.Path, ok = m["path"].(string); !ok || md.Path == "" {
		return md, fmt.Errorf("metadata.path is required")
	}

	// 必须有 operation (string → OperationType)
	opStr, ok := m["operation"].(string)
	if !ok || opStr == "" {
		return md, fmt.Errorf("metadata.operation is required")
	}
	md.Operation = protocols.OperationType(opStr)

	// 如果contentType为空，则设置为application/json
	if md.ContentType, ok = m["contentType"].(string); !ok || md.ContentType == "" {
		md.ContentType = "application/json"
	}

	if m["encryption"] != nil {
		md.Encryption = &protocols.Encryption{
			Algorithm:         m["encryption"].(map[string]interface{})["algorithm"].(string),
			SecretKey:         m["encryption"].(map[string]interface{})["secretKey"].(string),
			EncryptionPayload: m["encryption"].(map[string]interface{})["encryptionPayload"].(bool),
			EncryptionID:      m["encryption"].(map[string]interface{})["encryptionID"].(bool),
		}
	}

	return md, nil
}

func GetXid(username string, source string) (*protocols.XID[any], error) {
	// ctx := context.Background()
	// path := "/info"
	// if source != "" {
	// 	xidInfoRepository := xdb.NewXidInfoRepository()
	// 	path = fmt.Sprintf("/info/%s", source)
	// 	xid, err := xidInfoRepository.FindByName(ctx, username, path)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return xid, nil
	// }

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
