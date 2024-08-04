package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/launchdarkly-recruiting/re-coding-test-Simone-Caruso/db"
	"github.com/launchdarkly-recruiting/re-coding-test-Simone-Caruso/wsclient"
)

func TestMessageDispatch(t *testing.T) {
	jsonStr, _ := json.Marshal(&Score{1, "student", 0.4})
	tests := []struct {
		name     string
		payload  string
		result   []string
		hasError bool
	}{{
		"Single Message",
		fmt.Sprintf("\nevent: score\ndata: %s\n\n", jsonStr),
		[]string{"value"},
		false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(test.payload))
			c := wsclient.Client{}
			studentDB := db.New[string]()
			examDB := db.New[int]()
			d := &Dispatcher{studentDB, examDB}
			c.OnEvent(d.Listen)

			c.ReadStream(r)
			if len(studentDB.GetAll()) != 1 {
				t.Fatal("")
			}
			if _, a := studentDB.Get("student"); a != 0.4 {
				t.Fatal("")
			}
		})
	}

}
