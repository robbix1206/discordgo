package internal

import (
	"net/http"
	"time"

	"github.com/robbix1206/discordgo/logging"
)

// New create a new Client
func New(token string) *Client {
	return &Client{
		Ratelimiter:    NewRatelimiter(),
		MaxRestRetries: 3,
		Client:         &http.Client{Timeout: (20 * time.Second)},
		Token:          token,
	}
}

// GetLogLevel return the current log level of Client
func (s *Client) GetLogLevel() int {
	return s.LogLevel
}

func (s *Client) log(msgL int, format string, a ...interface{}) {
	logging.Log(s, msgL, format, a...)
}
