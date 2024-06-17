package serverStore

// ServerStore storage server data
type ServerStore interface {
	Hub() HubRepository
}
