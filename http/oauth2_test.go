package http

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestApplication(t *testing.T) {

	// Authentication Token pulled from environment variable DGU_TOKEN
	Token := os.Getenv("DGU_TOKEN")
	if Token == "" {
		t.Error("No token provided")
		return
	}

	// Create a new Discordgo session
	dg, err := discordgo.New(Token)
	if err != nil {
		t.Error(err)
		return
	}

	// Get a specific Application by it's ID
	_, err = dg.Application()
	if err != nil {
		t.Error(err)
	}
}
