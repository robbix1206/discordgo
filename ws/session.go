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

	// Exposed but should not be modified by User.

	// Whether the Data Websocket is ready
	DataReady bool // NOTE: Maye be deprecated soon

	// Status stores the currect status of the websocket connection
	// this is being tested, may stay, may go away.
	status int32

	// Whether the Voice Websocket is ready
	VoiceReady bool // NOTE: Deprecated.

	// Whether the UDP Connection is ready
	UDPReady bool // NOTE: Deprecated

	// Stores a mapping of guild id's to VoiceConnections
	VoiceConnections map[string]*VoiceConnection

	// Stores the last HeartbeatAck that was recieved (in UTC)
	LastHeartbeatAck time.Time

	// Stores the last Heartbeat sent (in UTC)
	LastHeartbeatSent time.Time

	// Event handlers
	handlersMu   sync.RWMutex
	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	// The websocket connection.
	wsConn *websocket.Conn

	// When nil, the session is not listening.
	listening chan interface{}

	// sequence tracks the current gateway api websocket sequence number
	sequence *int64

	// stores sessions current Discord Gateway
	gateway string

	// stores session ID of current Gateway connection
	sessionID string

	// used to make sure gateway websocket writes do not happen concurrently
	wsMutex sync.Mutex

	// userID is the ID of the current User
	userID string
}

func New(token string) *Session {
	return &Session{
		Compress:               true,
		ShouldReconnectOnError: true,
		ShardID:                0,
		ShardCount:             1,
		sequence:               new(int64),
		LastHeartbeatAck:       time.Now().UTC(),
	}
}

func (s *Session) GetLogLevel() int {
	return s.LogLevel
}

func (s *Session) log(msgL int, format string, a ...interface{}) {
	logging.Log(s, msgL, format, a)
}
