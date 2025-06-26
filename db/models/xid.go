package models

// Metadata {
// version: "1.0.0",
// CreatedAt: 1718630400,
// Operation: "create", // create, update, delete
//
//		Path: "/info/jumpserver",
//		ContentType: "application/json",
//

type Metadata struct {
	CreatedAt   int64  `json:"createdAt" bson:"createdAt"`
	Operation   string `json:"operation" bson:"operation"`
	Path        string `json:"path" bson:"path"`
	ContentType string `json:"contentType" bson:"contentType"`
}

type XID struct {
	Name     string      `json:"name" bson:"name"`
	Xid      string      `json:"xid" bson:"xid"`
	Version  string      `json:"version" bson:"version"`
	Metadata Metadata    `json:"metadata" bson:"metadata"`
	Payload  interface{} `json:"payload" bson:"payload"`
}
