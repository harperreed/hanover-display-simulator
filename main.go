package main

import (
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func main() {
	// Load configuration
	err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Initialize display
	initializeDisplay()

	go runWebServer()
	go processPackets()
	go readSerialPort()
	go testSimulator() // Run a test simulation

	// Keep the main goroutine running
	select {}
}
