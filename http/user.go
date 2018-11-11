package http

import (
	"bytes"
	"image"
	"net/url"
	"strconv"

	"github.com/robbix1206/discordgo/discord"
)

// User returns the user details of the given userID
// userID    : A user ID or "@me" which is a shortcut of current user ID
func (s *Session) User(userID string) (st *User, err error) {

	body, err := s.requestWithBucketID("GET", EndpointUser(userID), nil, EndpointUsers)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserUpdate updates a users settings.
func (s *Session) UserUpdate(username, avatar string) (st *discord.User, err error) {

	// NOTE: Avatar must be either the hash/id of existing Avatar or
	// data:image/png;base64,BASE64_STRING_OF_NEW_AVATAR_PNG
	// to set a new avatar.
	// If left blank, avatar will be set to null/blank

	data := struct {
		Username string `json:"username,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
	}{username, avatar}

	body, err := s.requestWithBucketID("PATCH", EndpointUser("@me"), data, EndpointUsers)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserGuilds returns an array of User structures for all guilds.
// limit     : The number guilds that can be returned. (max 100)
// beforeID  : If provided all guilds returned will be before given ID.
// afterID   : If provided all guilds returned will be after given ID.
// FIXME: change from UserGuild to Guild
func (s *Session) UserGuilds(limit int, beforeID, afterID string) (st []*Guild, err error) {

	v := url.Values{}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}

	uri := EndpointUserGuilds("@me")

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.requestWithBucketID("GET", uri, nil, EndpointUserGuilds(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildLeave leaves a Guild.
// guildID   : The ID of a Guild
func (s *Session) GuildLeave(guildID string) (err error) {
	_, err = s.requestWithBucketID("DELETE", EndpointUserGuild("@me", guildID), nil, EndpointUserGuild("", guildID))
	return
}

// UserChannels returns an array of Channel structures for all private
// channels.
func (s *Session) UserChannels() (st []*Channel, err error) {

	body, err := s.requestWithBucketID("GET", EndpointUserChannels("@me"), nil, EndpointUserChannels(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// recipientID : A user ID for the user to which this channel is opened with.
func (s *Session) UserChannelCreate(recipientID string) (st *Channel, err error) {

	data := struct {
		RecipientID string `json:"recipient_id"`
	}{recipientID}

	body, err := s.requestWithBucketID("POST", EndpointUserChannels("@me"), data, EndpointUserChannels(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// UserConnections returns the user's connections
func (s *Session) UserConnections() (conn []*UserConnection, err error) {
	response, err := s.requestWithBucketID("GET", EndpointUserConnections("@me"), nil, EndpointUserConnections("@me"))
	if err != nil {
		return nil, err
	}

	err = unmarshal(response, &conn)
	if err != nil {
		return
	}

	return
}

//TODO: Allow to create Grouped DM Channels

//TODO: Maybe separate CDN request as they are not part of the API

// UserAvatarDecode returns an image.Image of a user's Avatar
// user : The user which avatar should be retrieved
func (s *Session) UserAvatarDecode(u *discord.User) (img image.Image, err error) {
	body, err := s.requestWithBucketID("GET", EndpointUserAvatar(u.ID, u.Avatar), nil, EndpointUserAvatar("", ""))
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}
