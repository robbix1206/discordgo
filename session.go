package discordgo

import (
	"github.com/robbix1206/discordgo/http"
	"github.com/robbix1206/discordgo/ws"
)

//import "github.com/robbix1206/discordgo/ws"

type Session struct {
	HTTPSession *http.Session
	WsSession   *ws.Session
	State       *State
	// Stores a mapping of guild id's to VoiceConnections
	VoiceConnections map[string]*VoiceConnection
}

func (s *Session) Open() (err error) {
	gateway, err := s.HTTPSession.Gateway()
	//TODO: Open everything
	return err
}
