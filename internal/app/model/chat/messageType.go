package chat

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
