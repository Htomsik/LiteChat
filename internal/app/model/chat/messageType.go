package chat

import "encoding/json"

const (
	SystemUser = "System"
)

// MessageType Enum for type message
type MessageType int

const (
	TypeMessage MessageType = iota
	TypeUsersList
	TypeUserNameChanged
)

func (msgType MessageType) String() string {
	return [...]string{"Message", "UsersList", "UserNameChanged"}[msgType]
}

func (msgType MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(msgType.String())
}
