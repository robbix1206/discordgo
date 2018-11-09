package discord

// A VoiceState stores the voice states of Guilds
type VoiceState struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	//TODO: member
	SessionID string `json:"session_id"`
	Deaf      bool   `json:"deaf"`
	Mute      bool   `json:"mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	SelfMute  bool   `json:"self_mute"`
	Suppress  bool   `json:"suppress"`
}

// A VoiceRegion stores data for a specific voice region server.
type VoiceRegion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	//TODO: Remove hostname and port & add vip optimal deprecated and custom
	Hostname string `json:"sample_hostname"`
	Port     int    `json:"sample_port"`
}
