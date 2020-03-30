package main

import "github.com/lireza/lib/logging"

func main() {
	// Creating a logger with DEBUG log level.
	logger := logging.NewLogger(logging.DEBUG)

	// Must not be logged, because TRACE level is lower than DEBUG.
	logger.Trace("Calling remote service ...")

	logger.Debug("Request sent to %v", "192.168.1.10")
	logger.Info("Stopping the server ...")
}
