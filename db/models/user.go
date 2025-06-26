package models

import (
	"time"
)

// BaseUser 基础用户结构
type User struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Name      string    `json:"name" bson:"name"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
	Source    string    `json:"source" bson:"source"` // jumpserver, ldap, local, etc.
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// JumpServerUser JumpServer用户结构
type JumpServerUser struct {
	User               `bson:",inline"`
	Wechat             string   `json:"wechat" bson:"wechat"`
	Phone              string   `json:"phone" bson:"phone"`
	MFALevel           int      `json:"mfa_level" bson:"mfa_level"`
	SourceDisplay      string   `json:"source_display" bson:"source_display"`
	CanPublicKeyAuth   bool     `json:"can_public_key_auth" bson:"can_public_key_auth"`
	MFAEnabled         bool     `json:"mfa_enabled" bson:"mfa_enabled"`
	IsServiceAccount   bool     `json:"is_service_account" bson:"is_service_account"`
	IsValid            bool     `json:"is_valid" bson:"is_valid"`
	IsExpired          bool     `json:"is_expired" bson:"is_expired"`
	DateExpired        string   `json:"date_expired" bson:"date_expired"`
	DateJoined         string   `json:"date_joined" bson:"date_joined"`
	LastLogin          string   `json:"last_login" bson:"last_login"`
	CreatedBy          string   `json:"created_by" bson:"created_by"`
	Comment            string   `json:"comment" bson:"comment"`
	Groups             []string `json:"groups" bson:"groups"`
	GroupsDisplay      string   `json:"groups_display" bson:"groups_display"`
	SystemRoles        []string `json:"system_roles" bson:"system_roles"`
	OrgRoles           []string `json:"org_roles" bson:"org_roles"`
	SystemRolesDisplay string   `json:"system_roles_display" bson:"system_roles_display"`
	OrgRolesDisplay    string   `json:"org_roles_display" bson:"org_roles_display"`
	LoginBlocked       bool     `json:"login_blocked" bson:"login_blocked"`
}

// // LDAPUser LDAP用户结构
// type LDAPUser struct {
// 	BaseUser     `bson:",inline"`
// 	DN           string   `json:"dn" bson:"dn"`
// 	Department   string   `json:"department" bson:"department"`
// 	Title        string   `json:"title" bson:"title"`
// 	Manager      string   `json:"manager" bson:"manager"`
// 	Groups       []string `json:"groups" bson:"groups"`
// 	LastLogon    string   `json:"last_logon" bson:"last_logon"`
// 	PasswordAge  int      `json:"password_age" bson:"password_age"`
// 	AccountFlags int      `json:"account_flags" bson:"account_flags"`
// }

// // LocalUser 本地用户结构
// type LocalUser struct {
// 	BaseUser     `bson:",inline"`
// 	PasswordHash string    `json:"-" bson:"password_hash"` // 不在JSON中显示
// 	LastLogin    time.Time `json:"last_login" bson:"last_login"`
// 	LoginCount   int       `json:"login_count" bson:"login_count"`
// 	IsLocked     bool      `json:"is_locked" bson:"is_locked"`
// 	LockReason   string    `json:"lock_reason" bson:"lock_reason"`
// }
