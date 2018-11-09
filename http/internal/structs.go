// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains all structures for the discordgo package.  These
// may be moved about later into separate files but I find it easier to have
// them all located together.

package internal

import (
	"net/http"
	"sync"
	"time"
)

// Client is a structure cotaining the necessary to interact with discord HTTP API
type Client struct {
	sync.RWMutex

	// General configurable settings.

	// Authentication token for this session
	Token string

	// Debug for printing JSON request/responses
	LogLevel int

	// Exposed but should not be modified by User.

	// Max number of REST API retries
	MaxRestRetries int

	// The http client used for REST requests
	Client *http.Client

	// used to deal with rate limits
	Ratelimiter *RateLimiter
}

// RESTError stores error information about a request with a bad response code.
// Message is not always present, there are cases where api calls can fail
// without returning a json message.
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte

	Message *APIErrorMessage // Message may be nil.
}

// An APIErrorMessage is an api error message returned from discord
type APIErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// A TooManyRequests struct holds information received from Discord
// when receiving a HTTP 429 response.
type TooManyRequests struct {
	Bucket     string        `json:"bucket"`
	Message    string        `json:"message"`
	RetryAfter time.Duration `json:"retry_after"`
}
