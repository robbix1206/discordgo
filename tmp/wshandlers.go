
import "github.com/robbix1206/discordgo/logging"

// ------------------------------------------------------------------------------------------------
// Code related to voice connections that initiate over the data websocket
// ------------------------------------------------------------------------------------------------

// ChannelVoiceJoin joins the session user to a voice channel.
//
//    gID     : Guild ID of the channel to join.
//    cID     : Channel ID of the channel to join.
//    mute    : If true, you will be set to muted upon joining.
//    deaf    : If true, you will be set to deafened upon joining.
func (s *Session) ChannelVoiceJoin(gID, cID string, mute, deaf bool) (voice *VoiceConnection, err error) {

	s.log(logging.LogInformational, "called")

	s.RLock()
	voice, _ = s.VoiceConnections[gID]
	s.RUnlock()

	if voice == nil {
		voice = &VoiceConnection{}
		s.Lock()
		s.VoiceConnections[gID] = voice
		s.Unlock()
	}

	voice.Lock()
	voice.GuildID = gID
	voice.ChannelID = cID
	voice.deaf = deaf
	voice.mute = mute
	voice.session = s
	voice.Unlock()

	// Send the request to Discord that we want to join the voice channel
	data := voiceChannelJoinOp{4, voiceStateUpdateData{&gID, &cID, mute, deaf}}
	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(data)
	s.wsMutex.Unlock()
	if err != nil {
		return
	}

	// doesn't exactly work perfect yet.. TODO
	err = voice.waitUntilConnected()
	if err != nil {
		s.log(logging.LogWarning, "error waiting for voice to connect, %s", err)
		voice.Close()
		return
	}

	return
}

// onVoiceStateUpdate handles Voice State Update events on the data websocket.
func (s *Session) onVoiceStateUpdate(st *VoiceStateUpdate) {

	// If we don't have a connection for the channel, don't bother
	if st.ChannelID == "" {
		return
	}

	// Check if we have a voice connection to update
	s.RLock()
	voice, exists := s.VoiceConnections[st.GuildID]
	s.RUnlock()
	if !exists {
		return
	}

	// We only care about events that are about us.
	if s.userID != st.UserID {
		return
	}

	// Store the SessionID for later use.
	voice.Lock()
	voice.UserID = st.UserID
	voice.sessionID = st.SessionID
	voice.ChannelID = st.ChannelID
	voice.Unlock()
}

// onVoiceServerUpdate handles the Voice Server Update data websocket event.
//
// This is also fired if the Guild's voice region changes while connected
// to a voice channel.  In that case, need to re-establish connection to
// the new region endpoint.
func (s *Session) onVoiceServerUpdate(st *VoiceServerUpdate) {

	s.log(logging.LogInformational, "called")

	s.RLock()
	voice, exists := s.VoiceConnections[st.GuildID]
	s.RUnlock()

	// If no VoiceConnection exists, just skip this
	if !exists {
		return
	}

	// If currently connected to voice ws/udp, then disconnect.
	// Has no effect if not connected.
	voice.Close()

	// Store values for later use
	voice.Lock()
	voice.token = st.Token
	voice.endpoint = st.Endpoint
	voice.GuildID = st.GuildID
	voice.Unlock()

	// Open a connection to the voice server
	err := voice.open()
	if err != nil {
		s.log(logging.LogError, "onVoiceServerUpdate voice.open, %s", err)
	}
}

func (s *Session) onInterface(i interface{}) {
	switch t := i.(type) {
	case *Ready:
		/*
			for _, g := range t.Guilds {
				setGuildIds(g)
			}
		*/
		s.onReady(t)
		/*
			case *GuildCreate:
				setGuildIds(t.Guild)
			case *GuildUpdate:
				setGuildIds(t.Guild)
			case *VoiceServerUpdate:
				go s.onVoiceServerUpdate(t)
			case *VoiceStateUpdate:
				go s.onVoiceStateUpdate(t)
		*/
	}
	//FIXME: Permit to get all this events
	/*
		err := s.State.OnInterface(s, i)
		if err != nil {
			s.log(logging.LogDebug, "error dispatching internal event, %s", err)
		}
	*/
}

// setGuildIds will set the GuildID on all the members of a guild.
// This is done as event data does not have it set.
func setGuildIds(g *Guild) {
	for _, c := range g.Channels {
		c.GuildID = g.ID
	}
	// FIXME: Should not be useful
	/*
		for _, m := range g.Members {
			m.GuildID = g.ID
		}
	*/

	for _, vs := range g.VoiceStates {
		vs.GuildID = g.ID
	}
}

/* In Close
// I'm not sure if this is actually needed.
// if the gw reconnect works properly, voice should stay alive
// However, there seems to be cases where something "weird"
// happens.  So we're doing this for now just to improve
// stability in those edge cases.
s.RLock()
defer s.RUnlock()
for _, v := range s.VoiceConnections {
	s.log(logging.LogInformational, "reconnecting voice connection to guild %s", v.GuildID)
	go v.reconnect()
	// This is here just to prevent violently spamming the
	// voice reconnects
	time.Sleep(1 * time.Second)
}
*/
