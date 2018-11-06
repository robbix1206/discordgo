// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains high level helper functions and easy entry points for the
// entire discordgo package.  These functions are being developed and are very
// experimental at this point.  They will most likely change so please use the
// low level functions if that's a problem.

// Package discordgo provides Discord binding for Go
package discordgo

import (
	"github.com/robbix1206/discordgo/http"
)

// VERSION of DiscordGo, follows Semantic Versioning. (http://semver.org/)
const VERSION = "alpha"

// New creates a new Discord session with the given token
// All requests will use the token blindly,
// no verification of the token will be done and requests may fail.
// IF THE TOKEN IS FOR A BOT, IT MUST BE PREFIXED WITH `Bot `
// eg: `"Bot <token>"`
func New(token string) (s *Session) {
	// Create a Session interface.
	return &Session{
		//State:       NewState(),
		HTTPSession: http.New(token),
	}
}
