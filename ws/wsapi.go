// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains low level functions for interacting with the Discord
// data websocket interface.

package ws

import (
	"time"

	"github.com/robbix1206/discordgo/logging"
)

// HeartbeatLatency returns the latency between heartbeat acknowledgement and heartbeat send.
func (s *Socket) HeartbeatLatency() time.Duration {
	s.Lock()
	defer s.Unlock()
	return s.latency

}

func newUpdateStatusData(idle int, gameType ActivityType, game, url string) *StatusUpdateData {
	usd := &StatusUpdateData{
		Status: "online",
	}

	if idle > 0 {
		usd.IdleSince = &idle
	}

	if game != "" {
		usd.Game = &Activity{
			Name: game,
			Type: gameType,
			URL:  url,
		}
	}

	return usd
}

// UpdateStatus is used to update the user's status.
// If idle>0 then set status to idle.
// If game!="" then set game.
// if otherwise, set status to active, and no game.
func (s *Socket) UpdateStatus(idle int, game string) (err error) {
	return s.UpdateStatusComplex(*newUpdateStatusData(idle, ActivityTypeGame, game, ""))
}

// UpdateStreamingStatus is used to update the user's streaming status.
// If idle>0 then set status to idle.
// If game!="" then set game.
// If game!="" and url!="" then set the status type to streaming with the URL set.
// if otherwise, set status to active, and no game.
func (s *Socket) UpdateStreamingStatus(idle int, game string, url string) (err error) {
	gameType := ActivityTypeGame
	if url != "" {
		gameType = ActivityTypeStreaming
	}
	return s.UpdateStatusComplex(*newUpdateStatusData(idle, gameType, game, url))
}

// UpdateListeningStatus is used to set the user to "Listening to..."
// If game!="" then set to what user is listening to
// Else, set user to active and no game.
func (s *Socket) UpdateListeningStatus(game string) (err error) {
	return s.UpdateStatusComplex(*newUpdateStatusData(0, ActivityTypeListening, game, ""))
}

// UpdateStatusComplex allows for sending the raw status update data untouched by discordgo.
func (s *Socket) UpdateStatusComplex(usd StatusUpdateData) (err error) {
	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}

	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(gatewayPayload{Op: 3, Data: usd})
	s.wsMutex.Unlock()

	return
}

// RequestGuildMembers requests guild members from the gateway
// The gateway responds with GuildMembersChunk events
// guildID  : The ID of the guild to request members of
// query    : String that username starts with, leave empty to return all members
// limit    : Max number of items to return, or 0 to request all members matched
func (s *Socket) RequestGuildMembers(guildID, query string, limit int) (err error) {
	s.log(logging.LogInformational, "called")

	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}

	data := requestGuildMembersData{
		GuildID: guildID,
		Query:   query,
		Limit:   limit,
	}

	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(gatewayPayload{Op: 8, Data: data})
	s.wsMutex.Unlock()

	return
}

// ChannelVoiceJoinManual initiates a voice session to a voice channel, but does not complete it.
//
// This should only be used when the VoiceServerUpdate will be intercepted and used elsewhere.
//
//    gID     : Guild ID of the channel to join.
//    cID     : Channel ID of the channel to join.
//    mute    : If true, you will be set to muted upon joining.
//    deaf    : If true, you will be set to deafened upon joining.
func (s *Socket) ChannelVoiceJoinManual(gID, cID string, mute, deaf bool) (err error) {
	s.log(logging.LogInformational, "called")

	s.RLock()
	defer s.RUnlock()
	if s.wsConn == nil {
		return ErrWSNotFound
	}
	// Send the request to Discord that we want to join the voice channel
	data := voiceStateUpdateData{&gID, &cID, mute, deaf}
	s.wsMutex.Lock()
	err = s.wsConn.WriteJSON(gatewayPayload{Op: 4, Data: data})
	s.wsMutex.Unlock()
	if err != nil {
		return
	}
	return
}
