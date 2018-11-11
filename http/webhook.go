package http

import "github.com/robbix1206/discordgo/discord"

// WebhookCreate returns a new Webhook.
// channelID: The ID of a Channel.
// name     : The name of the webhook.
// avatar   : The avatar of the webhook.
func (s *Session) WebhookCreate(channelID, name, avatar string) (st *Webhook, err error) {

	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	body, err := s.requestWithBucketID("POST", EndpointChannelWebhooks(channelID), data, EndpointChannelWebhooks(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// ChannelWebhooks returns all webhooks for a given channel.
// channelID: The ID of a channel.
func (s *Session) ChannelWebhooks(channelID string) (st []*Webhook, err error) {

	body, err := s.requestWithBucketID("GET", EndpointChannelWebhooks(channelID), nil, EndpointChannelWebhooks(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// GuildWebhooks returns all webhooks for a given guild.
// guildID: The ID of a Guild.
func (s *Session) GuildWebhooks(guildID string) (st []*Webhook, err error) {

	body, err := s.requestWithBucketID("GET", EndpointGuildWebhooks(guildID), nil, EndpointGuildWebhooks(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// Webhook returns a webhook for a given ID
// webhookID: The ID of a webhook.
func (s *Session) Webhook(webhookID string) (st *Webhook, err error) {

	body, err := s.requestWithBucketID("GET", EndpointWebhook(webhookID), nil, EndpointWebhooks)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookWithToken returns a webhook for a given ID
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook.
func (s *Session) WebhookWithToken(webhookID, token string) (st *Webhook, err error) {

	body, err := s.requestWithBucketID("GET", EndpointWebhookToken(webhookID, token), nil, EndpointWebhookToken("", ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookEdit updates an existing Webhook.
// webhookID: The ID of a webhook.
// name     : The name of the webhook.
// avatar   : The avatar of the webhook.
func (s *Session) WebhookEdit(webhookID, name, avatar, channelID string) (st *discord.User, err error) {

	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}{name, avatar, channelID}

	body, err := s.requestWithBucketID("PATCH", EndpointWebhook(webhookID), data, EndpointWebhooks)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookEditWithToken updates an existing Webhook with an auth token.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook.
// name     : The name of the webhook.
// avatar   : The avatar of the webhook.
func (s *Session) WebhookEditWithToken(webhookID, token, name, avatar string) (err error) {

	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	_, err = s.requestWithBucketID("PATCH", EndpointWebhookToken(webhookID, token), data, EndpointWebhookToken("", ""))

	return
}

// WebhookDelete deletes a webhook for a given ID
// webhookID: The ID of a webhook.
func (s *Session) WebhookDelete(webhookID string) (err error) {

	_, err = s.requestWithBucketID("DELETE", EndpointWebhook(webhookID), nil, EndpointWebhooks)

	return
}

// WebhookDeleteWithToken deletes a webhook for a given ID with an auth token.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook.
func (s *Session) WebhookDeleteWithToken(webhookID, token string) (st *Webhook, err error) {

	body, err := s.requestWithBucketID("DELETE", EndpointWebhookToken(webhookID, token), nil, EndpointWebhookToken("", ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

// WebhookExecute executes a webhook.
// webhookID: The ID of a webhook.
// token    : The auth token for the webhook
func (s *Session) WebhookExecute(webhookID, token string, wait bool, data *WebhookParams) (err error) {
	uri := EndpointWebhookToken(webhookID, token)

	if wait {
		uri += "?wait=true"
	}

	_, err = s.requestWithBucketID("POST", uri, data, EndpointWebhookToken("", ""))

	return
}
