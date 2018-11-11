package discordgo

import (
	"github.com/robbix1206/discordgo/http"
	"github.com/robbix1206/discordgo/ws"
)

// Session represent a simple session (A HTTP and a websocket) to communicate with discord
type Session struct {
	HTTPSession *http.Session
	WsSession   *ws.Socket
	State       *State
	// Stores a mapping of guild id's to VoiceConnections
	VoiceConnections map[string]*VoiceConnection
}
