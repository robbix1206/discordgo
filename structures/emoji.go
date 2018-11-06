package structures

// Emoji struct holds data related to Emoji's
type Emoji struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
	//TODO: Add user
	RequireColons bool `json:"require_colons"`
	Managed       bool `json:"managed"`
	Animated      bool `json:"animated"`
}
