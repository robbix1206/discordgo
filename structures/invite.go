package structures

// A Invite stores all data related to a specific Discord Guild or Channel invite.
type Invite struct {
	Code    string   `json:"code"`
	Guild   *Guild   `json:"guild"`
	Channel *Channel `json:"channel"`

	// Will only be filled when using InviteWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	// Metadata part
	Inviter   *User     `json:"inviter"`
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt Timestamp `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}
