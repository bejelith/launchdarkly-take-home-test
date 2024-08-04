package main

import (
	"encoding/json"
	log "log/slog"

	"github.com/launchdarkly-recruiting/re-coding-test-Simone-Caruso/db"
	"github.com/launchdarkly-recruiting/re-coding-test-Simone-Caruso/wsclient"
)

// Score is the event wire format specified in README.md
type Score struct {
	Exam      int     `json:"exam"`
	StudentId string  `json:"studentId"`
	Score     float64 `json:"score"`
}

// Dispatcher will deliver all received Messages to required backend DBs
type Dispatcher struct {
	StudentDB *db.Average[string]
	ExamDB    *db.Average[int]
}

// Listen is the receiver method for new Messages to dispatch to the underlying DBs
func (d *Dispatcher) Listen(m *wsclient.Message) {
	score := &Score{}
	err := json.Unmarshal([]byte(m.Data), score)
	if err != nil {
		log.Error("Failed to parse json message", "error", err)
		return
	}
	d.StudentDB.Add(score.StudentId, score.Score)
	d.ExamDB.Add(score.Exam, score.Score)
}
