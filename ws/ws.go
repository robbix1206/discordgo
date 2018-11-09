package ws

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robbix1206/discordgo/discord"
	"github.com/robbix1206/discordgo/logging"
)

// ErrWSAlreadyOpen is thrown when you attempt to open
// a websocket that already is open.
var ErrWSAlreadyOpen = errors.New("web socket already opened")

// ErrWSNotFound is thrown when you attempt to use a websocket
// that doesn't exist
var ErrWSNotFound = errors.New("no websocket connection exists")

// ErrWSShardBounds is thrown when you try to use a shard ID that is
// less than the total shard count
var ErrWSShardBounds = errors.New("ShardID must be less than ShardCount")

// routineSystem get Evend sent by Discord and sends regular heartbeats
// to Discord so it knows the client is still connected.
// If you do not send these heartbeats Discord will
// disconnect the websocket connection after a few seconds.
func (s *Session) routineSystem(wsConn *websocket.Conn, heartbeatIntervalMsec time.Duration) {

	s.log(logging.LogInformational, "called")

	if wsConn == nil {
		return
	}

	eventChan := make(chan *Event)
	go s.eventListener(wsConn, eventChan)

	ticker := time.NewTicker(heartbeatIntervalMsec * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Lock()
			waitingAck := s.waitingAck
			s.waitingAck = true
			s.Unlock()
			if waitingAck {
				s.log(logging.LogError, "no ACK between heartbeats, triggering a reconnection")
				s.autoReconnect()
				return
			}
			if s.heartbeat() != nil {
				s.log(logging.LogError, "wasn't able to heartbeat")
			}
		case e, ok := <-eventChan:
			if !ok {
				s.RLock()
				sameConnection := s.wsConn == wsConn
				s.RUnlock()
				if sameConnection {
					s.log(logging.LogWarning, "error reading from gateway %s websocket, Reconnect", s.gateway)
					s.Reconnect()
				}
				return
			}
			s.onEvent(e)
		}
	}
}

// listen polls the websocket connection for events, it will stop when the
// listening channel is closed, or an error occurs.
func (s *Session) eventListener(wsConn *websocket.Conn, eventChan chan<- *Event) {

	s.log(logging.LogInformational, "called")
	defer close(eventChan)

	for {
		e, err := s.getEvent()
		if err != nil {
			s.RLock()
			s.log(logging.LogWarning, "error reading from gateway %s websocket, %s", s.gateway, err)
			s.RUnlock()
			return
		}
		eventChan <- e
	}
}

func (s *Session) heartbeat() (err error) {
	sequence := atomic.LoadInt64(&s.sequence)
	s.log(logging.LogDebug, "sending gateway websocket heartbeat seq %d", sequence)
	s.Lock()
	s.lastHeartbeatSent = time.Now().UTC()
	s.Unlock()
	s.wsMutex.Lock()
	if sequence != 0 {
		err = s.wsConn.WriteJSON(gatewayPayload{Op: 1, Data: sequence})
	} else {
		err = s.wsConn.WriteJSON(gatewayPayload{Op: 1, Data: nil})
	}
	s.wsMutex.Unlock()
	if err != nil {
		s.log(logging.LogError, "error sending heartbeat to gateway %s, %s", s.gateway, err)
	}
	return err
}

// identify sends the identify packet to the gateway
func (s *Session) identify() error {
	data := identifyData{
		Token: s.Token,
		Properties: identifyConnectionProperties{
			OS:      runtime.GOOS,
			Browser: "Discordgo v" + VERSION,
			Device:  "Discordgo v" + VERSION,
		},
		Compress:       s.Compress,
		LargeThreshold: 250,
		Shard:          nil,
	}
	if s.ShardCount > 1 {
		if s.ShardID >= s.ShardCount {
			return ErrWSShardBounds
		}
		data.Shard = &[2]int{s.ShardID, s.ShardCount}
	}

	s.wsMutex.Lock()
	err := s.wsConn.WriteJSON(gatewayPayload{
		Op:   2,
		Data: data,
	})
	s.wsMutex.Unlock()

	return err
}

func (s *Session) resume() error {
	// Send Op 6 Resume Packet
	data := resumeData{
		Token:     s.Token,
		SessionID: s.sessionID,
		Sequence:  atomic.LoadInt64(&s.sequence),
	}

	s.log(logging.LogInformational, "sending resume packet to gateway")
	s.wsMutex.Lock()
	err := s.wsConn.WriteJSON(gatewayPayload{Op: 6, Data: data})
	s.wsMutex.Unlock()
	if err != nil {
		err = fmt.Errorf("error sending gateway resume packet, %s, %s", s.gateway, err)
	}
	return err
}

func (s *Session) getEvent() (e *Event, err error) {
	messageType, message, err := s.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var reader io.Reader
	reader = bytes.NewBuffer(message)

	// If this is a compressed message, uncompress it.
	if messageType == websocket.BinaryMessage {
		z, err := zlib.NewReader(reader)
		if err != nil {
			s.log(logging.LogError, "error uncompressing websocket message, %s", err)
			return nil, err
		}
		reader = z

		defer func() {
			err := z.Close()
			if err != nil {
				s.log(logging.LogWarning, "error closing zlib, %s", err)
			}
		}()

	}
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&e); err != nil {
		s.log(logging.LogError, "error decoding websocket message, %s", err)
		return nil, err
	}
	s.log(logging.LogDebug, "Op: %d, Seq: %d, Type: %s, Data: %s\n\n", e.Operation, e.Sequence, e.Type, string(e.RawData))
	return e, err

}

// onEvent is the "event handler" for all messages received on the
// Discord Gateway API websocket connection.
//
// If you use the AddHandler() function to register a handler for a
// specific event this function will pass the event along to that handler.
//
// If you use the AddHandler() function to register a handler for the
// "OnEvent" event then all events will be passed to that handler.

func (s *Session) onEvent(e *Event) (err error) {
	switch e.Operation {
	case 0:
		atomic.StoreInt64(&s.sequence, e.Sequence)
	case 1:
		s.log(logging.LogInformational, "sending heartbeat in response to Op1")
		err = s.heartbeat()
	case 7:
		s.log(logging.LogInformational, "Reconnecting in response to Op7")
		s.Reconnect()
	case 9:
		// Managed directly in resumeSystem
	case 10:
		// Managed directly in handshake
	case 11:
		s.Lock()
		s.latency = time.Now().UTC().Sub(s.lastHeartbeatSent)
		s.waitingAck = false
		s.Unlock()
		s.log(logging.LogDebug, "got heartbeat ACK")
	default:
		s.log(logging.LogWarning, "unknown Op: %d, Seq: %d, Type: %s, Data: %s", e.Operation, e.Sequence, e.Type, string(e.RawData))
	}
	return err
}

func (s *Session) dispatch(e *Event) error {

	// Do not try to Dispatch a non-Dispatch Message
	if e.Operation != 0 {
		return fmt.Errorf("This is not a dispatch event Op: %d", e.Operation)
	}

	// Map event to registered event handlers and pass it along to any registered handlers.
	if eh, ok := registeredInterfaceProviders[e.Type]; ok {
		e.Struct = eh.New()

		// Attempt to unmarshal our event.
		if err := json.Unmarshal(e.RawData, e.Struct); err != nil {
			s.log(logging.LogError, "error unmarshalling %s event, %s", e.Type, err)
			return err
		}
		s.handleEvent(e.Type, e.Struct)
	} else {
		s.log(logging.LogWarning, "unknown event: Op: %d, Seq: %d, Type: %s, Data: %s", e.Operation, e.Sequence, e.Type, string(e.RawData))
	}

	// For legacy reasons, we send the raw event also, this could be useful for handling unknown events.
	s.handleEvent(eventEventType, e)

	return nil
}

func (s *Session) identifySystem() error {
	// Send Op 2 Identity Packet
	err := s.identify()
	if err != nil {
		err = fmt.Errorf("error sending identify packet to gateway, %s, %s", s.gateway, err)
		return err
	}
	for {
		e, err := s.getEvent()
		if err != nil {
			return err
		}
		if e.Operation == 0 {
			if e.Type != `READY` {
				s.log(logging.LogError, "Expected READY, instead got:\n%#v\n", e)
				return fmt.Errorf("Unexpected Dispatch Received")
			}
			s.onEvent(e)
			return nil
		} else if e.Operation == 9 {
			// Is dispatch when oyou hit the ratelimit of 2 identify in less than 5 seconds
			return fmt.Errorf("Invalid Session")
		} else if e.Operation == 1 || e.Operation == 11 {
			s.onEvent(e)
		} else {
			s.log(logging.LogError, "Expected READY, instead got:\n%#v\n", e)
			return fmt.Errorf("Unexpected Packet received")
		}
	}
}

func (s *Session) resumeSystem() error {
	err := s.resume()
	if err != nil {
		return err
	}
	for {
		e, err := s.getEvent()
		if err != nil {
			return err
		}
		if e.Operation == 0 {
			if e.Type != `RESUME` {
				s.log(logging.LogError, "Expected RESUMED, instead got:\n%#v\n", e)
				return fmt.Errorf("Unexpected Dispatch Received")
			}
			s.onEvent(e)
			return nil
		} else if e.Operation == 9 {
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
			s.log(logging.LogInformational, "sending identify packet to gateway in response to Op9")
			return s.identifySystem()
		} else if e.Operation == 1 || e.Operation == 11 {
			s.onEvent(e)
		} else {
			s.log(logging.LogError, "Expected RESUMED, instead got:\n%#v\n", e)
			return fmt.Errorf("Unexpected Packet received")
		}
	}
}

func (s *Session) getHello(e *Event) (h *helloData, err error) {
	if e.Operation != 10 {
		return nil, fmt.Errorf("expecting Op 10, got Op %d instead", e.Operation)
	}
	h = new(helloData)
	s.log(logging.LogInformational, "Op 10 Hello Packet received from Discord")
	if err := json.Unmarshal(e.RawData, h); err != nil {
		err = fmt.Errorf("error unmarshalling helloData, %s", err)
		return nil, err
	}
	return h, err
}

func (s *Session) handshake() error {
	// The first response from Discord should be an Op 10 (Hello) Packet.
	e, err := s.getEvent()
	if err != nil {
		return err
	}
	h, err := s.getHello(e)
	if err != nil {
		return err
	}
	// Now we send either an Op 2 Identity if this is a brand new
	// connection or Op 6 Resume if we are resuming an existing connection.
	if s.sessionID == "" {
		err = s.identifySystem()
	} else {
		err = s.resumeSystem()
	}
	if err != nil {
		return err
	}
	s.log(logging.LogInformational, "We are now connected to Discord, emitting connect event")
	s.handleEvent(connectEventType, &Connect{})

	// Start sending heartbeats and reading messages from Discord.
	go s.routineSystem(s.wsConn, h.HeartbeatInterval)

	s.log(logging.LogInformational, "exiting")
	return nil
}

// Open creates a websocket connection to Discord.
// See: https://discordapp.com/developers/docs/topics/gateway#connecting
func (s *Session) Open() (err error) {
	s.log(logging.LogInformational, "called")
	// Prevent Open or other major Session functions from
	// being called while Open is still running.
	s.Lock()
	defer s.Unlock()

	// If the websock is already open, bail out here.
	if s.wsConn != nil {
		return ErrWSAlreadyOpen
	}

	if s.gateway == "" {
		url, err := s.getGateway()
		if err != nil {
			return err
		}
		s.gateway += url + "?v=" + discord.APIVersion + "&encoding=json"
	}

	// Connect to the Gateway
	s.log(logging.LogInformational, "connecting to gateway %s", s.gateway)
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	s.wsConn, _, err = websocket.DefaultDialer.Dial(s.gateway, header)
	if err != nil {
		s.log(logging.LogWarning, "error connecting to gateway %s, %s", s.gateway, err)
		s.gateway = "" // clear cached gateway
		s.wsConn = nil // Just to be safe.
		return err
	}

	s.wsConn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	err = s.handshake()
	if err != nil {
		s.wsConn.Close()
		s.wsConn = nil
	}
	return err
}

// Close closes the websocket
func (s *Session) Close() (err error) {

	s.log(logging.LogInformational, "called")
	s.Lock()

	if s.wsConn != nil {
		s.log(logging.LogInformational, "sending close frame")
		// To cleanly close a connection, a client should send a close
		// frame and wait for the server to close the connection.
		s.wsMutex.Lock()
		err := s.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		s.wsMutex.Unlock()
		if err != nil {
			s.log(logging.LogInformational, "error closing websocket, %s", err)
		}
		// TODO: Wait for Discord to actually close the connection.
		time.Sleep(1 * time.Second)

		s.log(logging.LogInformational, "closing gateway websocket")
		err = s.wsConn.Close()
		if err != nil {
			s.log(logging.LogInformational, "error closing websocket, %s", err)
		}

		s.wsConn = nil
	}

	s.Unlock()

	s.log(logging.LogInformational, "emit disconnect event")
	s.handleEvent(disconnectEventType, &Disconnect{})

	return
}

func (s *Session) Reconnect() {

	s.log(logging.LogInformational, "called")

	err := s.Close()
	if err != nil {
		s.log(logging.LogWarning, "error closing session connection, %s", err)
	}

	wait := time.Duration(1)

	for {
		s.log(logging.LogInformational, "trying to reconnect to gateway")
		//TODO: Set a maximum number of try
		err = s.Open()
		if err == nil {
			s.log(logging.LogInformational, "successfully reconnected to gateway")
			return
		}
		s.log(logging.LogError, "error reconnecting to gateway, %s", err)

		<-time.After(wait * time.Second)
		wait *= 2
		if wait > 600 {
			wait = 600
		}
	}
}

func (s *Session) autoReconnect() {
	if s.ShouldReconnectOnError {
		s.Reconnect()
	}
}
