// Package wsclient implements the server-side event protocol
package wsclient

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	log "log/slog"
)

const eventField = "event: "
const dataField = "data: "

// Client is a Server-Sent compliant client, will keep reading from a Stream until EOF or CLose() is called
type Client struct {
	started   int32
	listeners []func(*Message)
}

// Close will terminate the stream
func (c *Client) Close() {
	atomic.StoreInt32(&c.started, 0)
}

// OnEvent will add a listener to all received events
func (c *Client) OnEvent(f func(*Message)) {
	c.listeners = append(c.listeners, f)
}

// ReadStream reads Server-Sent events from a Stream, can be run in a separate goroutine
func (c *Client) ReadStream(r io.Reader) {
	atomic.StoreInt32(&c.started, 1)
	scanner := bufio.NewScanner(r)
	buffer := []string{}
	for scanner.Scan() && atomic.LoadInt32(&c.started) == 1 {
		line := scanner.Text()
		if len(line) == 0 { // New message delimiter has been read
			if len(buffer) > 0 {
				message, err := readMessage(buffer)
				if err != nil {
					log.Error("failed to parse message", "error", err)
				} else {
					for _, l := range c.listeners {
						l(message)
					}
				}
				buffer = []string{}
			}
		} else {
			buffer = append(buffer, line)
		}
	}
}

// Message represents a Server-Sent event, MessageType is optional but not for this specific application
type Message struct {
	MessageType string
	Data        string
}

// readMessage parses a multiline Server-Set payload and returns a Message
func readMessage(m []string) (*Message, error) {
	// We need at least one valid `event` field (this is specific to this exercise as we always expect a event field)
	if len(m) < 2 {
		return nil, fmt.Errorf("Incomplete message received")
	}
	if !strings.HasPrefix(m[0], eventField) {
		return nil, fmt.Errorf("Missing event type field in message")
	}
	message := &Message{}
	message.MessageType = strings.Trim(m[0][len(eventField):], " ")

	payload := bytes.NewBufferString("")
	for i := 1; i < len(m); i++ {
		line := m[i]
		if !strings.HasPrefix(line, dataField) {
			return nil, fmt.Errorf("Invalid message fields found")
		}
		payload.WriteString(line[len(dataField):])
	}
	message.Data = payload.String()
	return message, nil
}
