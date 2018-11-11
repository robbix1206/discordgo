package http

import (
	"bytes"
	"fmt"
	"image"
	"net/url"
	"strconv"

	"github.com/robbix1206/discordgo/discord"
)

// A GuildParams stores all the data needed to update discord guild settings
type GuildParams struct {
	Name                        string             `json:"name,omitempty"`
	Region                      string             `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	AfkChannelID                string             `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                `json:"afk_timeout,omitempty"`
	Icon                        string             `json:"icon,omitempty"`
	OwnerID                     string             `json:"owner_id,omitempty"`
	Splash                      string             `json:"splash,omitempty"`
}

// A GuildEmbed stores data for a guild embed.
type GuildEmbed struct {
	Enabled   bool   `json:"enabled"`
	ChannelID string `json:"channel_id"`
}

// A GuildAuditLog stores data for a guild audit log.
type GuildAuditLog struct {
	Webhooks []struct {
		ChannelID string `json:"channel_id"`
		GuildID   string `json:"guild_id"`
		ID        string `json:"id"`
		Avatar    string `json:"avatar"`
		Name      string `json:"name"`
	} `json:"webhooks,omitempty"`
	Users []struct {
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
		ID            string `json:"id"`
		Avatar        string `json:"avatar"`
	} `json:"users,omitempty"`
	AuditLogEntries []struct {
		TargetID string `json:"target_id"`
		Changes  []struct {
			NewValue interface{} `json:"new_value"`
			OldValue interface{} `json:"old_value"`
			Key      string      `json:"key"`
		} `json:"changes,omitempty"`
		UserID     string `json:"user_id"`
		ID         string `json:"id"`
		ActionType int    `json:"action_type"`
		Options    struct {
			DeleteMembersDay string `json:"delete_member_days"`
			MembersRemoved   string `json:"members_removed"`
			ChannelID        string `json:"channel_id"`
			Count            string `json:"count"`
			ID               string `json:"id"`
			Type             string `json:"type"`
			RoleName         string `json:"role_name"`
		} `json:"options,omitempty"`
		Reason string `json:"reason"`
	} `json:"audit_log_entries"`
}

// Block contains Discord Audit Log Action Types
const (
	AuditLogActionGuildUpdate = 1

	AuditLogActionChannelCreate          = 10
	AuditLogActionChannelUpdate          = 11
	AuditLogActionChannelDelete          = 12
	AuditLogActionChannelOverwriteCreate = 13
	AuditLogActionChannelOverwriteUpdate = 14
	AuditLogActionChannelOverwriteDelete = 15

	AuditLogActionMemberKick       = 20
	AuditLogActionMemberPrune      = 21
	AuditLogActionMemberBanAdd     = 22
	AuditLogActionMemberBanRemove  = 23
	AuditLogActionMemberUpdate     = 24
	AuditLogActionMemberRoleUpdate = 25

	AuditLogActionRoleCreate = 30
	AuditLogActionRoleUpdate = 31
	AuditLogActionRoleDelete = 32

	AuditLogActionInviteCreate = 40
	AuditLogActionInviteUpdate = 41
	AuditLogActionInviteDelete = 42

	AuditLogActionWebhookCreate = 50
	AuditLogActionWebhookUpdate = 51
	AuditLogActionWebhookDelete = 52

	AuditLogActionEmojiCreate = 60
	AuditLogActionEmojiUpdate = 61
	AuditLogActionEmojiDelete = 62

	AuditLogActionMessageDelete = 72
)

// Guild returns a Guild structure of a specific Guild.
// guildID   : The ID of a Guild
func (s *Session) Guild(guildID string) (st *Guild, err error) {
	body, err := s.requestWithBucketID("GET", EndpointGuild(guildID), nil, EndpointGuild(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildCreate creates a new Guild
// name      : A name for the Guild (2-100 characters)
func (s *Session) GuildCreate(name string) (st *Guild, err error) {

	data := struct {
		Name string `json:"name"`
	}{name}

	body, err := s.requestWithBucketID("POST", EndpointGuildCreate, data, EndpointGuildCreate)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildEdit edits a new Guild
// guildID   : The ID of a Guild
// g 		 : A GuildParams struct with the values Name, Region and VerificationLevel defined.
func (s *Session) GuildEdit(guildID string, g GuildParams) (st *Guild, err error) {

	// Bounds checking for VerificationLevel, interval: [0, 3]
	if g.VerificationLevel != nil {
		val := *g.VerificationLevel
		if val < 0 || val > 3 {
			err = ErrVerificationLevelBounds
			return
		}
	}

	//Bounds checking for regions
	if g.Region != "" {
		isValid := false
		regions, _ := s.VoiceRegions()
		for _, r := range regions {
			if g.Region == r.ID {
				isValid = true
			}
		}
		if !isValid {
			var valid []string
			for _, r := range regions {
				valid = append(valid, r.ID)
			}
			err = fmt.Errorf("Region not a valid region (%q)", valid)
			return
		}
	}

	body, err := s.requestWithBucketID("PATCH", EndpointGuild(guildID), g, EndpointGuild(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildDelete deletes a Guild.
// guildID   : The ID of a Guild
func (s *Session) GuildDelete(guildID string) (st *Guild, err error) {

	body, err := s.requestWithBucketID("DELETE", EndpointGuild(guildID), nil, EndpointGuild(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildBans returns an array of User structures for all bans of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildBans(guildID string) (st []*GuildBan, err error) {

	body, err := s.requestWithBucketID("GET", EndpointGuildBans(guildID), nil, EndpointGuildBans(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildBanCreate bans the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
// days      : The number of days of previous comments to delete.
func (s *Session) GuildBanCreate(guildID, userID string, days int) (err error) {
	return s.GuildBanCreateWithReason(guildID, userID, "", days)
}

// GuildBanCreateWithReason bans the given user from the given guild also providing a reaso.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
// reason    : The reason for this ban
// days      : The number of days of previous comments to delete.
func (s *Session) GuildBanCreateWithReason(guildID, userID, reason string, days int) (err error) {

	uri := EndpointGuildBan(guildID, userID)

	queryParams := url.Values{}
	if days > 0 {
		queryParams.Set("delete-message-days", strconv.Itoa(days))
	}
	if reason != "" {
		queryParams.Set("reason", reason)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	_, err = s.requestWithBucketID("PUT", uri, nil, EndpointGuildBan(guildID, ""))
	return
}

// GuildBanDelete removes the given user from the guild bans
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildBanDelete(guildID, userID string) (err error) {

	_, err = s.requestWithBucketID("DELETE", EndpointGuildBan(guildID, userID), nil, EndpointGuildBan(guildID, ""))
	return
}

// GuildMembers returns a list of members for a guild.
//  guildID  : The ID of a Guild.
//  after    : The id of the member to return members after
//  limit    : max number of members to return (max 1000)
func (s *Session) GuildMembers(guildID string, after string, limit int) (st []*Member, err error) {

	uri := EndpointGuildMembers(guildID)

	v := url.Values{}

	if after != "" {
		v.Set("after", after)
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.requestWithBucketID("GET", uri, nil, EndpointGuildMembers(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildMember returns a member of a guild.
//  guildID   : The ID of a Guild.
//  userID    : The ID of a User
func (s *Session) GuildMember(guildID, userID string) (st *Member, err error) {

	body, err := s.requestWithBucketID("GET", EndpointGuildMember(guildID, userID), nil, EndpointGuildMember(guildID, ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildMemberAdd force joins a user to the guild.
//  accessToken   : Valid access_token for the user.
//  guildID       : The ID of a Guild.
//  userID        : The ID of a User.
//  nick          : Value to set users nickname to
//  roles         : A list of role ID's to set on the member.
//  mute          : If the user is muted.
//  deaf          : If the user is deafened.
func (s *Session) GuildMemberAdd(accessToken, guildID, userID, nick string, roles []string, mute, deaf bool) (err error) {

	data := struct {
		AccessToken string   `json:"access_token"`
		Nick        string   `json:"nick,omitempty"`
		Roles       []string `json:"roles,omitempty"`
		Mute        bool     `json:"mute,omitempty"`
		Deaf        bool     `json:"deaf,omitempty"`
	}{accessToken, nick, roles, mute, deaf}

	_, err = s.requestWithBucketID("PUT", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	if err != nil {
		return err
	}

	return err
}

// GuildMemberDelete removes the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
func (s *Session) GuildMemberDelete(guildID, userID string) (err error) {

	return s.GuildMemberDeleteWithReason(guildID, userID, "")
}

// GuildMemberDeleteWithReason removes the given user from the given guild.
// guildID   : The ID of a Guild.
// userID    : The ID of a User
// reason    : The reason for the kick
func (s *Session) GuildMemberDeleteWithReason(guildID, userID, reason string) (err error) {

	uri := EndpointGuildMember(guildID, userID)
	if reason != "" {
		uri += "?reason=" + url.QueryEscape(reason)
	}

	_, err = s.requestWithBucketID("DELETE", uri, nil, EndpointGuildMember(guildID, ""))
	return
}

// GuildMemberEdit edits the roles of a member.
// guildID  : The ID of a Guild.
// userID   : The ID of a User.
// roles    : A list of role ID's to set on the member.
func (s *Session) GuildMemberEdit(guildID, userID string, roles []string) (err error) {

	data := struct {
		Roles []string `json:"roles"`
	}{roles}

	_, err = s.requestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	if err != nil {
		return
	}

	return
}

// GuildMemberMove moves a guild member from one voice channel to another/none
//  guildID   : The ID of a Guild.
//  userID    : The ID of a User.
//  channelID : The ID of a channel to move user to, or null?
// NOTE : I am not entirely set on the name of this function and it may change
// prior to the final 1.0.0 release of Discordgo
func (s *Session) GuildMemberMove(guildID, userID, channelID string) (err error) {

	data := struct {
		ChannelID string `json:"channel_id"`
	}{channelID}

	_, err = s.requestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	if err != nil {
		return
	}

	return
}

// GuildMemberNickname updates the nickname of a guild member
// guildID   : The ID of a guild
// userID    : The ID of a user
// userID    : The ID of a user or "@me" which is a shortcut of the current user ID
func (s *Session) GuildMemberNickname(guildID, userID, nickname string) (err error) {

	data := struct {
		Nick string `json:"nick"`
	}{nickname}

	if userID == "@me" {
		userID += "/nick"
	}

	_, err = s.requestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	return
}

// GuildMemberRoleAdd adds the specified role to a given member
//  guildID   : The ID of a Guild.
//  userID    : The ID of a User.
//  roleID 	  : The ID of a Role to be assigned to the user.
func (s *Session) GuildMemberRoleAdd(guildID, userID, roleID string) (err error) {

	_, err = s.requestWithBucketID("PUT", EndpointGuildMemberRole(guildID, userID, roleID), nil, EndpointGuildMemberRole(guildID, "", ""))

	return
}

// GuildMemberRoleRemove removes the specified role to a given member
//  guildID   : The ID of a Guild.
//  userID    : The ID of a User.
//  roleID 	  : The ID of a Role to be removed from the user.
func (s *Session) GuildMemberRoleRemove(guildID, userID, roleID string) (err error) {

	_, err = s.requestWithBucketID("DELETE", EndpointGuildMemberRole(guildID, userID, roleID), nil, EndpointGuildMemberRole(guildID, "", ""))

	return
}

// GuildChannels returns an array of Channel structures for all channels of a
// given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildChannels(guildID string) (st []*Channel, err error) {

	body, err := s.request("GET", EndpointGuildChannels(guildID), "", nil, EndpointGuildChannels(guildID), 0)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildChannelCreateData is provided to GuildChannelCreateComplex
type GuildChannelCreateData struct {
	Name                 string                 `json:"name"`
	Type                 ChannelType            `json:"type"`
	Topic                string                 `json:"topic,omitempty"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
}

// GuildChannelCreateComplex creates a new channel in the given guild
// guildID      : The ID of a Guild
// data         : A data struct describing the new Channel, Name and Type are mandatory, other fields depending on the type
func (s *Session) GuildChannelCreateComplex(guildID string, data GuildChannelCreateData) (st *Channel, err error) {
	body, err := s.requestWithBucketID("POST", EndpointGuildChannels(guildID), data, EndpointGuildChannels(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildChannelCreate creates a new channel in the given guild
// guildID   : The ID of a Guild.
// name      : Name of the channel (2-100 chars length)
// ctype     : Type of the channel
func (s *Session) GuildChannelCreate(guildID, name string, ctype ChannelType) (st *Channel, err error) {
	return s.GuildChannelCreateComplex(guildID, GuildChannelCreateData{
		Name: name,
		Type: ctype,
	})
}

// GuildChannelsReorder updates the order of channels in a guild
// guildID   : The ID of a Guild.
// channels  : Updated channels.
func (s *Session) GuildChannelsReorder(guildID string, channels []*Channel) (err error) {

	data := make([]struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	_, err = s.requestWithBucketID("PATCH", EndpointGuildChannels(guildID), data, EndpointGuildChannels(guildID))
	return
}

// GuildInvites returns an array of Invite structures for the given guild
// guildID   : The ID of a Guild.
func (s *Session) GuildInvites(guildID string) (st []*Invite, err error) {
	body, err := s.requestWithBucketID("GET", EndpointGuildInvites(guildID), nil, EndpointGuildInvites(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildRoles returns all roles for a given guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildRoles(guildID string) (st []*Role, err error) {

	body, err := s.requestWithBucketID("GET", EndpointGuildRoles(guildID), nil, EndpointGuildRoles(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return // TODO return pointer
}

// GuildRoleCreate returns a new Guild Role.
// guildID: The ID of a Guild.
func (s *Session) GuildRoleCreate(guildID string) (st *Role, err error) {

	body, err := s.requestWithBucketID("POST", EndpointGuildRoles(guildID), nil, EndpointGuildRoles(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildRoleEdit updates an existing Guild Role with new values
// guildID   : The ID of a Guild.
// roleID    : The ID of a Role.
// name      : The name of the Role.
// color     : The color of the role (decimal, not hex).
// hoist     : Whether to display the role's users separately.
// perm      : The permissions for the role.
// mention   : Whether this role is mentionable
func (s *Session) GuildRoleEdit(guildID, roleID, name string, color int, hoist bool, perm int, mention bool) (st *Role, err error) {

	// Prevent sending a color int that is too big.
	if color > 0xFFFFFF {
		err = fmt.Errorf("color value cannot be larger than 0xFFFFFF")
		return nil, err
	}

	data := struct {
		Name        string `json:"name"`        // The role's name (overwrites existing)
		Color       int    `json:"color"`       // The color the role should have (as a decimal, not hex)
		Hoist       bool   `json:"hoist"`       // Whether to display the role's users separately
		Permissions int    `json:"permissions"` // The overall permissions number of the role (overwrites existing)
		Mentionable bool   `json:"mentionable"` // Whether this role is mentionable
	}{name, color, hoist, perm, mention}

	body, err := s.requestWithBucketID("PATCH", EndpointGuildRole(guildID, roleID), data, EndpointGuildRole(guildID, ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildRoleReorder reoders guild roles
// guildID   : The ID of a Guild.
// roles     : A list of ordered roles.
func (s *Session) GuildRoleReorder(guildID string, roles []*Role) (st []*Role, err error) {

	body, err := s.requestWithBucketID("PATCH", EndpointGuildRoles(guildID), roles, EndpointGuildRoles(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildRoleDelete deletes an existing role.
// guildID   : The ID of a Guild.
// roleID    : The ID of a Role.
func (s *Session) GuildRoleDelete(guildID, roleID string) (err error) {

	_, err = s.requestWithBucketID("DELETE", EndpointGuildRole(guildID, roleID), nil, EndpointGuildRole(guildID, ""))

	return
}

// GuildPruneCount Returns the number of members that would be removed in a prune operation.
// Requires 'KICK_MEMBER' permission.
// guildID	: The ID of a Guild.
// days		: The number of days to count prune for (1 or more).
func (s *Session) GuildPruneCount(guildID string, days uint32) (count uint32, err error) {
	count = 0

	if days <= 0 {
		err = ErrPruneDaysBounds
		return
	}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	uri := EndpointGuildPrune(guildID) + "?days=" + strconv.FormatUint(uint64(days), 10)
	body, err := s.requestWithBucketID("GET", uri, nil, EndpointGuildPrune(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &p)
	if err != nil {
		return
	}

	count = p.Pruned

	return
}

// GuildPrune Begin as prune operation. Requires the 'KICK_MEMBERS' permission.
// Returns an object with one 'pruned' key indicating the number of members that were removed in the prune operation.
// guildID	: The ID of a Guild.
// days		: The number of days to count prune for (1 or more).
func (s *Session) GuildPrune(guildID string, days uint32) (count uint32, err error) {

	count = 0

	if days <= 0 {
		err = ErrPruneDaysBounds
		return
	}

	data := struct {
		days uint32
	}{days}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	body, err := s.requestWithBucketID("POST", EndpointGuildPrune(guildID), data, EndpointGuildPrune(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &p)
	if err != nil {
		return
	}

	count = p.Pruned

	return
}

// GuildIntegrations returns an array of Integrations for a guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildIntegrations(guildID string) (st []*discord.Integration, err error) {

	body, err := s.requestWithBucketID("GET", EndpointGuildIntegrations(guildID), nil, EndpointGuildIntegrations(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildIntegrationCreate creates a Guild Integration.
// guildID          : The ID of a Guild.
// integrationType  : The Integration type.
// integrationID    : The ID of an integration.
func (s *Session) GuildIntegrationCreate(guildID, integrationType, integrationID string) (err error) {

	data := struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{integrationType, integrationID}

	_, err = s.requestWithBucketID("POST", EndpointGuildIntegrations(guildID), data, EndpointGuildIntegrations(guildID))
	return
}

// GuildIntegrationEdit edits a Guild Integration.
// guildID              : The ID of a Guild.
// integrationType      : The Integration type.
// integrationID        : The ID of an integration.
// expireBehavior	      : The behavior when an integration subscription lapses (see the integration object documentation).
// expireGracePeriod    : Period (in seconds) where the integration will ignore lapsed subscriptions.
// enableEmoticons	    : Whether emoticons should be synced for this integration (twitch only currently).
func (s *Session) GuildIntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool) (err error) {

	data := struct {
		ExpireBehavior    int  `json:"expire_behavior"`
		ExpireGracePeriod int  `json:"expire_grace_period"`
		EnableEmoticons   bool `json:"enable_emoticons"`
	}{expireBehavior, expireGracePeriod, enableEmoticons}

	_, err = s.requestWithBucketID("PATCH", EndpointGuildIntegration(guildID, integrationID), data, EndpointGuildIntegration(guildID, ""))
	return
}

// GuildIntegrationDelete removes the given integration from the Guild.
// guildID          : The ID of a Guild.
// integrationID    : The ID of an integration.
func (s *Session) GuildIntegrationDelete(guildID, integrationID string) (err error) {

	_, err = s.requestWithBucketID("DELETE", EndpointGuildIntegration(guildID, integrationID), nil, EndpointGuildIntegration(guildID, ""))
	return
}

// GuildIntegrationSync syncs an integration.
// guildID          : The ID of a Guild.
// integrationID    : The ID of an integration.
func (s *Session) GuildIntegrationSync(guildID, integrationID string) (err error) {

	_, err = s.requestWithBucketID("POST", EndpointGuildIntegrationSync(guildID, integrationID), nil, EndpointGuildIntegration(guildID, ""))
	return
}

// GuildIcon returns an image.Image of a guild icon.
// guildID   : The ID of a Guild.
func (s *Session) GuildIcon(guildID string) (img image.Image, err error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return
	}

	if g.Icon == "" {
		err = ErrGuildNoIcon
		return
	}

	body, err := s.requestWithBucketID("GET", EndpointGuildIcon(guildID, g.Icon), nil, EndpointGuildIcon(guildID, ""))
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

// GuildSplash returns an image.Image of a guild splash image.
// guildID   : The ID of a Guild.
func (s *Session) GuildSplash(guildID string) (img image.Image, err error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return
	}

	if g.Splash == "" {
		err = ErrGuildNoSplash
		return
	}

	body, err := s.requestWithBucketID("GET", EndpointGuildSplash(guildID, g.Splash), nil, EndpointGuildSplash(guildID, ""))
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

// GuildEmbed returns the embed for a Guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildEmbed(guildID string) (st *GuildEmbed, err error) {

	body, err := s.requestWithBucketID("GET", EndpointGuildEmbed(guildID), nil, EndpointGuildEmbed(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildEmbedEdit returns the embed for a Guild.
// guildID   : The ID of a Guild.
func (s *Session) GuildEmbedEdit(guildID string, enabled bool, channelID string) (err error) {

	data := GuildEmbed{enabled, channelID}

	_, err = s.requestWithBucketID("PATCH", EndpointGuildEmbed(guildID), data, EndpointGuildEmbed(guildID))
	return
}

// GuildAuditLog returns the audit log for a Guild.
// guildID     : The ID of a Guild.
// userID      : If provided the log will be filtered for the given ID.
// beforeID    : If provided all log entries returned will be before the given ID.
// actionType  : If provided the log will be filtered for the given Action Type.
// limit       : The number messages that can be returned. (default 50, min 1, max 100)
func (s *Session) GuildAuditLog(guildID, userID, beforeID string, actionType, limit int) (st *GuildAuditLog, err error) {

	uri := EndpointGuildAuditLogs(guildID)

	v := url.Values{}
	if userID != "" {
		v.Set("user_id", userID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if actionType > 0 {
		v.Set("action_type", strconv.Itoa(actionType))
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if len(v) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, v.Encode())
	}

	body, err := s.requestWithBucketID("GET", uri, nil, EndpointGuildAuditLogs(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// GuildEmojiCreate creates a new emoji
// guildID : The ID of a Guild.
// name    : The Name of the Emoji.
// image   : The base64 encoded emoji image, has to be smaller than 256KB.
// roles   : The roles for which this emoji will be whitelisted, can be nil.
func (s *Session) GuildEmojiCreate(guildID, name, image string, roles []string) (emoji *Emoji, err error) {

	data := struct {
		Name  string   `json:"name"`
		Image string   `json:"image"`
		Roles []string `json:"roles,omitempty"`
	}{name, image, roles}

	body, err := s.requestWithBucketID("POST", EndpointGuildEmojis(guildID), data, EndpointGuildEmojis(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &emoji)
	return
}

// GuildEmojiEdit modifies an emoji
// guildID : The ID of a Guild.
// emojiID : The ID of an Emoji.
// name    : The Name of the Emoji.
// roles   : The roles for which this emoji will be whitelisted, can be nil.
func (s *Session) GuildEmojiEdit(guildID, emojiID, name string, roles []string) (emoji *Emoji, err error) {

	data := struct {
		Name  string   `json:"name"`
		Roles []string `json:"roles,omitempty"`
	}{name, roles}

	body, err := s.requestWithBucketID("PATCH", EndpointGuildEmoji(guildID, emojiID), data, EndpointGuildEmojis(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &emoji)
	return
}

// GuildEmojiDelete deletes an Emoji.
// guildID : The ID of a Guild.
// emojiID : The ID of an Emoji.
func (s *Session) GuildEmojiDelete(guildID, emojiID string) (err error) {

	_, err = s.requestWithBucketID("DELETE", EndpointGuildEmoji(guildID, emojiID), nil, EndpointGuildEmojis(guildID))
	return
}
