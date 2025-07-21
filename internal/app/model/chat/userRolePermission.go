package chat

import "encoding/json"

type RolePermission int

const (
	PermissionNone RolePermission = iota
	PermissionSendMessage
	PermissionDeleteChat
)

func (permType RolePermission) String() string {
	return [...]string{"None", "SendMessage", "DeleteChat"}[permType]
}

func (permType RolePermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(permType.String())
}
