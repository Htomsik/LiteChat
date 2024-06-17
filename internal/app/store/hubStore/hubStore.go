package hubStore

// HubStore storage hub data
type HubStore interface {
	Client() ClientRepository
}
