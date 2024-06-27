package chat

import "encoding/json"

type RolePermission int

const (
	PermissionSendMessage RolePermission = iota
	PermissionDeleteChat
)

func (permType RolePermission) String() string {
	return [...]string{"SendMessage", "DeleteChat"}[permType]
}

func (permType RolePermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(permType.String())
}
