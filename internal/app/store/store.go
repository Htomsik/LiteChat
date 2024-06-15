package store

// Store storage server data
type Store interface {
	Hub() HubRepository
}
