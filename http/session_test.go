package http

import (
	"fmt"
	"os"
)

var (
	dg    *Session // Stores a global discordgo user session
	dgBot *Session // Stores a global discordgo bot session

	envToken    = os.Getenv("DGU_TOKEN")  // Token to use when authenticating the user account
	envBotToken = os.Getenv("DGB_TOKEN")  // Token to use when authenticating the bot account
	envGuild    = os.Getenv("DG_GUILD")   // Guild ID to use for tests
	envChannel  = os.Getenv("DG_CHANNEL") // Channel ID to use for tests
	envAdmin    = os.Getenv("DG_ADMIN")   // User ID of admin user to use for tests
)

func init() {
	fmt.Println("Init is being called.")
	if envBotToken != "" {
		dgBot = New(envBotToken)
	}
	dg = New(envToken)
}
