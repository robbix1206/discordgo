package structures

// A PresenceUpdate stores the online, offline, or idle and game status of Guild members.
type PresenceUpdate struct {
	User       *User     `json:"user"`
	Roles      []string  `json:"roles"`
	Game       *Activity `json:"game"`
	GuildID    string    `json:"guild_id"`
	Status     Status    `json:"status"`
	Activities *Activity `json:"activities"`
}

// Status type definition
type Status string

// Constants for Status with the different current available status
const (
	StatusOnline       Status = "online"
	StatusIdle         Status = "idle"
	StatusDoNotDisturb Status = "dnd"
	StatusInvisible    Status = "invisible"
	StatusOffline      Status = "offline"
)

// A Activity struct holds the name of the "playing .." game for a user
type Activity struct {
	Name          string             `json:"name"`
	Type          ActivityType       `json:"type"`
	URL           string             `json:"url,omitempty"`
	TimeStamps    ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID string             `json:"application_id,omitempty"`
	Details       string             `json:"details,omitempty"`
	State         string             `json:"state,omitempty"`
	Party         ActivityParty      `json:"party,omitempty"`
	Assets        ActivityAssets     `json:"assets,omitempty"`
	Secrets       ActivitySecret     `json:"secrets,omitempty"`
	Instance      bool               `json:"instance,omitempty"`
	Flags         ActivityFlag       `json:"flags,omitempty"`
}

// ActivityType is the type of "game" (see GameType* consts) in the Game struct
type ActivityType int

// Valid ActivityType values
const (
	ActivityTypeGame ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
)

// A ActivityTimestamps struct contains start and end times used in the rich presence "playing .." Game
type ActivityTimestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

// A ActivityParty struct contains information about the party
type ActivityParty struct {
	ID string `json:"id,omitempty"`
	//First element is current size second element is max size
	Size [2]int `json:"size,omitempty"`
}

// An ActivityAssets struct contains assets and labels used in the rich presence "playing .." Game
type ActivityAssets struct {
	LargeImageID string `json:"large_image,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}

// A ActivitySecret struct contains information activity secrets
type ActivitySecret struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

// ActivityFlag type definition
type ActivityFlag int

// Constants for ActivityFlag
const (
	ActivityFlagInstance ActivityFlag = 1 << iota
	ActivityFlagJoin
	ActivityFlagSpectate
	ActivityFlagJoinRequest
	ActivityFlagSync
	ActivityFlagPlay
)

// GatewayBotResponse stores the data for the gateway/bot response
type GatewayBotResponse struct {
	URL    string `json:"url"`
	Shards int    `json:"shards"`
	//TODO: Add session_start_limit
}
