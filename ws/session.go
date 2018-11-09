package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robbix1206/discordgo/logging"
)

type Session struct {
	sync.RWMutex

	// Authentication token for this session
	Token string

	// General configurable settings.
	LogLevel int

	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool

	// Should the session request compressed websocket data.
	Compress bool

	// Sharding
	ShardID    int
	ShardCount int

	// Whether or not to call event handlers synchronously.
	// e.g false = launch event handlers in their own goroutines.
	SyncEvents bool

	// Stores the Duration between an heartbeat and it's ACK
	latency time.Duration

	// Stores the last Heartbeat sent (in UTC)
	lastHeartbeatSent time.Time

	// ReceivedHeartbeatAck check if we received an ACK between two heatbeats
	waitingAck bool

	// Event handlers
	handlersMu   sync.RWMutex
	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	// The websocket connection.
	wsConn *websocket.Conn

	// sequence tracks the current gateway api websocket sequence number
	sequence int64

	// stores session ID of current Gateway connection
	sessionID string

	// used to make sure gateway websocket writes do not happen concurrently
	wsMutex sync.Mutex

	// userID is the ID of the current User
	userID string

	// stores sessions current cached Discord Gateway
	gateway string

	// getGateway is a function that allow to get the gateway URL to connect
	getGateway func() (string, error)
}

func New(token string) *Session {
	return &Session{
		Compress:               true,
		ShouldReconnectOnError: true,
		ShardID:                0,
		ShardCount:             1,
	}
}

func (s *Session) GetLogLevel() int {
	return s.LogLevel
}

func (s *Session) log(msgL int, format string, a ...interface{}) {
	logging.Log(s, msgL, format, a...)
}
