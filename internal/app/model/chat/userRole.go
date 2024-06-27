package chat

// UserRole permissions type of users
type UserRole struct {
	Name        string           `toml:"name"`
	IsAdmin     bool             `toml:"isAdmin"`
	Permissions []RolePermission `toml:"permissions"`
}
