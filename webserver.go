package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	templates *template.Template
	clients   = make(map[chan string]bool)
	clientsMutex sync.Mutex
)

func init() {
	templates = template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/display.html",
	))
}

func runWebServer() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		display.mu.Lock()
		defer display.mu.Unlock()

		err := templates.ExecuteTemplate(c.Writer, "layout.html", display.pixels)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error executing template")
			return
		}
	})

	r.GET("/events", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		clientChan := make(chan string)
		clientsMutex.Lock()
		clients[clientChan] = true
		clientsMutex.Unlock()

		defer func() {
			clientsMutex.Lock()
			delete(clients, clientChan)
			clientsMutex.Unlock()
		}()

		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	// ... (keep other routes like /packets and /display)

	go func() {
		for range time.Tick(100 * time.Millisecond) {
			updateClients()
		}
	}()

	if err := r.Run(config.WebPort); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

func updateClients() {
	display.mu.Lock()
	defer display.mu.Unlock()

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "display", display.pixels)
	if err != nil {
		log.Errorf("Error executing template: %v", err)
		return
	}
	displayHTML := buf.String()

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for clientChan := range clients {
		select {
		case clientChan <- displayHTML:
		default:
			// Client is not ready to receive, skip this update
		}
	}
}

func notifyNewPacket() {
	go updateClients()
}
