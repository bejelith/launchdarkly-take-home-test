package main

import (
	"flag"
	log "log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/launchdarkly-recruiting/re-coding-test-Simone-Caruso/db"
	"github.com/launchdarkly-recruiting/re-coding-test-Simone-Caruso/wsclient"
)

var sourceURL = flag.String("url", "https://live-test-scores.herokuapp.com/scores", "URL to stream events from")
var listenAddr = flag.String("listen", "localhost:8080", "Listen address")
var debugLogs = flag.Bool("debug", false, "Enable debug log")

// HTTPResponse is used to return samples and their average
type HTTPResponse struct {
	Avg     float64
	Samples []float64
}

func main() {
	flag.Parse()

	// Setup basic logging
	if *debugLogs {
		log.SetLogLoggerLevel(log.LevelDebug)
		gin.SetMode(gin.DebugMode)
	} else {
		log.SetLogLoggerLevel(log.LevelInfo)
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	resp, err := http.Get(*sourceURL)
	if err != nil {
		log.Error("Failed to open stream", "url", *sourceURL, "error", err)
		os.Exit(1)
	}

	// Initialize WS client
	wsc := wsclient.Client{}

	// Setup graceful shutdown
	setupSignalHandler(func(s os.Signal) {
		log.Info("Signal received, stopping server", "signal", s)
		wsc.Close()
		resp.Body.Close()
		os.Exit(0)
	})

	// Setup core application, DBs and Dispatcher
	studentDB := db.New[string]()
	examDB := db.New[int]()
	d := &Dispatcher{studentDB, examDB}
	wsc.OnEvent(d.Listen)
	go wsc.ReadStream(resp.Body)

	// Init REST apis
	router.GET("/students", getStudentsHandler(studentDB))
	router.GET("/students/:id", getStudentHandler(studentDB))
	router.GET("/exams", getExamsHandler(examDB))
	router.GET("/exams/:id", getExamHandler(examDB))

	if err := router.Run(*listenAddr); err != nil {
		log.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

// Builder methods for HTTP handlers
func getStudentsHandler(db *db.Average[string]) func(*gin.Context) {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, db.GetAll())
	}
}

func getStudentHandler(db *db.Average[string]) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		samples, avg := db.Get(id)
		c.IndentedJSON(http.StatusOK, HTTPResponse{avg, samples})
	}
}

func getExamHandler(db *db.Average[int]) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		numericID, err := strconv.Atoi(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Id needs to be of numeric type")
			return
		}
		samples, avg := db.Get(numericID)
		c.IndentedJSON(http.StatusOK, HTTPResponse{avg, samples})
	}
}

func getExamsHandler(db *db.Average[int]) func(*gin.Context) {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, db.GetAll())
	}
}
