package proxy

import (
	"fmt"

	"github.com/chatmcp/mcprouter/service/mcpclient"
)

// SSESession is a session for SSE request
type SSESession struct {
	writer   *SSEWriter
	done     chan struct{} // done channel
	messages chan string   // event queue
	key      string        // client request key
	command  string
	client   *mcpclient.StdioClient
}

// NewSSESession will create a new SSE session
func NewSSESession(w *SSEWriter, key string, command string) *SSESession {
	return &SSESession{
		writer:   w,
		done:     make(chan struct{}),
		messages: make(chan string, 100), // store messages
		key:      key,
		command:  command,
		client:   nil,
	}
}

// Key returns the key of the session
func (s *SSESession) Key() string {
	return s.key
}

// Command returns the command of the session
func (s *SSESession) Command() string {
	return s.command
}

// SetClient sets the client of the session
func (s *SSESession) SetClient(client *mcpclient.StdioClient) {
	s.client = client
}

// Client returns the client of the session
func (s *SSESession) Client() *mcpclient.StdioClient {
	return s.client
}

// Messages returns the messages channel of the session
func (s *SSESession) Messages() chan string {
	return s.messages
}

// SendMessage sends a message to the session
func (s *SSESession) SendMessage(message string) {
	select {
	case s.messages <- message:
		// send message to client ok
	case <-s.done:
		fmt.Printf("session is closed\n")
	default:
		fmt.Printf("message channel is full\n")
	}
}

// Close closes the session
func (s *SSESession) Close() {
	s.CloseClient()
	close(s.done)
}

func (s *SSESession) CloseClient() {
	if s.client != nil {
		s.client.Close()
	}
}

// Done returns the done channel of the session
func (s *SSESession) Done() chan struct{} {
	return s.done
}
