package http

import "github.com/robbix1206/discordgo/discord"

// An Application struct stores values for a Discord OAuth2 Application
type Application struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	Icon                string        `json:"icon,omitempty"`
	Description         string        `json:"description"`
	RPCOrigins          *[]string     `json:"rpc_origins"`
	BotPublic           bool          `json:"bot_public"`
	BotRequireCodeGrant bool          `json:"bot_require_code_grant"`
	Owner               *discord.User `json:"owner"`
}

// Application returns an Application structure of this Application
func (s *Session) Application() (st *Application, err error) {
	body, err := s.requestWithBucketID("GET", EndpointApplication("@me"), nil, EndpointApplication(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}
