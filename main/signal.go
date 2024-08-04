package main

import (
	"os"
	"os/signal"
	"syscall"
)

// setupSignalHandler will use f to handle os Signals
func setupSignalHandler(f func(os.Signal)) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		f(s)
	}()
}
