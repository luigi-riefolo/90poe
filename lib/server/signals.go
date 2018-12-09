package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/common/log"
)

// HandleSignals listens for a list of system signal,
// if any is emitted then fn is executed in order
// to gently shut down the system.
func HandleSignals(fn func()) {

	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGSTOP,
		syscall.SIGABRT,
		syscall.SIGQUIT)

	go func() {
		for sig := range ch {
			log.Infof("Caught signal '%v', shutting down", sig)

			// tidy up
			fn()

			// that's all folks
			os.Exit(1)
		}
	}()
}
