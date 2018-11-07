package ws

import "github.com/robbix1206/discordgo/structures"

type Guild = structures.Guild
type Activity = structures.Activity
type ActivityType = structures.ActivityType

const (
	ActivityTypeGame      = structures.ActivityTypeGame
	ActivityTypeStreaming = structures.ActivityTypeStreaming
	ActivityTypeListening = structures.ActivityTypeListening
)
