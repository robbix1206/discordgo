package discord

// A User stores all data for an individual Discord user.
type User struct {
	// The ID of the user.
	ID string `json:"id"`

	// The user's username.
	Username string `json:"username"`

	// The discriminator of the user (4 numbers after name).
	Discriminator string `json:"discriminator"`

	// The hash of the user's avatar. Use Session.UserAvatar
	// to retrieve the avatar itself.
	Avatar string `json:"avatar"`

	// Whether the user is a bot.
	Bot bool `json:"bot"`

	// Whether the user has multi-factor authentication enabled.
	MFAEnabled bool `json:"mfa_enabled"`

	// The user's chosen language option.
	Locale string `json:"locale"`

	// Whether the user's email is verified.
	Verified bool `json:"verified"`

	// The email of the user. This is only present when
	// the application possesses the email scope for the user.
	Email string `json:"email"`
}

// UserConnection is a Connection returned from the UserConnections endpoint
type UserConnection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
}
