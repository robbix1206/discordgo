package http

import (
	"github.com/robbix1206/discordgo/http/internal"
)

type Session struct {
	internal *internal.Client
}

func New(token string) *Session {
	return &Session{
		internal: internal.New(token),
	}
}

func unmarshal(data []byte, v interface{}) error {
	return internal.Unmarshal(data, v)
}

func (s *Session) RequestWithBucketID(method, urlStr string, data interface{}, bucketID string) (response []byte, err error) {
	return s.internal.RequestWithBucketID(method, urlStr, data, bucketID)
}

func (s *Session) request(method, urlStr, contentType string, b []byte, bucketID string, sequence int) (response []byte, err error) {
	return s.internal.Request(method, urlStr, contentType, b, bucketID, sequence)
}
