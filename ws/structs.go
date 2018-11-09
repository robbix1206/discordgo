package ws

import (
	"time"

	"github.com/robbix1206/discordgo/discord"
)

type Guild = discord.Guild
type Activity = discord.Activity
type ActivityType = discord.ActivityType

const (
	ActivityTypeGame      = discord.ActivityTypeGame
	ActivityTypeStreaming = discord.ActivityTypeStreaming
	ActivityTypeListening = discord.ActivityTypeListening
)

type gatewayPayload struct {
	Op        int         `json:"op"`
	Data      interface{} `json:"d"`
	Sequence  int         `json:"s,omitempty"`
	EventName string      `json:"t,omitempty"`
}

type helloData struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

type identifyConnectionProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type identifyData struct {
	Token          string                       `json:"token"`
	Properties     identifyConnectionProperties `json:"properties"`
	Compress       bool                         `json:"compress"`
	LargeThreshold int                          `json:"large_threshold"`
	Shard          *[2]int                      `json:"shard,omitempty"`
}

// StatusUpdateData ia provided to UpdateStatusComplex()
type StatusUpdateData struct {
	IdleSince *int      `json:"since"`
	Game      *Activity `json:"game"`
	AFK       bool      `json:"afk"`
	Status    string    `json:"status"`
}

type resumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Sequence  int64  `json:"seq"`
}

type voiceStateUpdateData struct {
	GuildID   *string `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
}

type requestGuildMembersData struct {
	GuildID string `json:"guild_id"`
	Query   string `json:"query"`
	Limit   int    `json:"limit"`
}
