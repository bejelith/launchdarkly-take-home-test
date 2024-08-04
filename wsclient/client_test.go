package wsclient

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadStream(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		result   []string
		hasError bool
	}{{
		"Single Message",
		"\nevent: score\ndata: value\n\n",
		[]string{"value"},
		false,
	}, {
		"Two Messages",
		"\nevent: score\ndata: 1\n\n\nevent: score\ndata: 2\n\n",
		[]string{"1", "2"},
		false,
	}, {
		"Multiple Data Fields",
		"\nevent: score\ndata: 1\ndata: 2\n\n",
		[]string{"12"},
		false,
	}, {
		"Missing event type",
		"\ndata: 1\ndata: 2\n\n",
		[]string{},
		false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(test.payload))
			c := &Client{}
			messages := []*Message{}
			c.OnEvent(func(m *Message) {
				messages = append(messages, m)
			})
			c.ReadStream(r)
			if len(messages) != len(test.result) {
				t.Fatal("Some messages where lost")
			}
			for i, m := range messages {
				if m.Data != test.result[i] {
					t.Fatalf("Expected %s but got %s", test.result[i], m.Data)
				}
			}
		})
	}

}

func TestReadMessage(t *testing.T) {
	tests := []struct {
		name     string
		payload  []string
		result   string
		hasError bool
	}{{
		"Valid Message",
		[]string{"event: score", "data: 1", "data: 2"},
		"12",
		false,
	}, {
		"No Data Message",
		[]string{"event: score"},
		"",
		true,
	}, {
		"Invalid data",
		[]string{"event: score", "Invalid data"},
		"",
		true,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, e := readMessage(test.payload)
			if test.hasError {
				if e == nil || m != nil {
					t.Fatal(e)
				}
			} else if e != nil {
				t.Fatal(e)
			} else {
				if m.Data != test.result {
					fmt.Printf("expected %s but got %s", test.result, m.Data)
				}
			}

		})
	}
}
