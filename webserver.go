package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	templates    *template.Template
	clients      = make(map[chan string]bool)
	clientsMutex sync.Mutex
)

func init() {
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"iterate": func(count int) []int {
			var i int
			var Items []int
			for i = 0; i < (count); i++ {
				Items = append(Items, i)
			}
			return Items
		},
	}).ParseFiles(
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

			jsonData := pixelsToJSON(display.pixels)

			err := templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
				"Pixels":   display.pixels,
				"Rows":     config.Rows,
				"Columns":  config.Columns,
				"JSONData": jsonData,
			})
			if err != nil {
				log.Errorf("Error executing template: %v", err)
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
				close(clientChan)
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

	r.GET("/packets", func(c *gin.Context) {
		packetInfos := make([]struct {
			Timestamp time.Time
			Length    int
		}, len(packetLog))
		for i, p := range packetLog {
			packetInfos[i] = struct {
				Timestamp time.Time
				Length    int
			}{
				Timestamp: p.Timestamp,
				Length:    len(p.Data),
			}
		}
		c.JSON(http.StatusOK, packetInfos)
	})

	r.GET("/display", func(c *gin.Context) {
		display.mu.Lock()
		defer display.mu.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"pixels": display.pixels,
			"json":   pixelsToJSON(display.pixels),
		})
	})

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

	jsonData := pixelsToJSON(display.pixels)

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "display", gin.H{
		"Pixels":  display.pixels,
		"Rows":    config.Rows,
		"Columns": config.Columns,
	})
	if err != nil {
		log.Errorf("Error executing template: %v", err)
		return
	}
	displayHTML := buf.String()

	updateData := struct {
		HTML string `json:"html"`
		JSON string `json:"json"`
	}{
		HTML: displayHTML,
		JSON: jsonData,
	}

	updateJSON, err := json.Marshal(updateData)
	if err != nil {
		log.Errorf("Error marshaling update data: %v", err)
		return
	}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for clientChan := range clients {
		select {
		case clientChan <- string(updateJSON):
			log.Debug("Sent update to client")
		default:
			log.Debug("Client not ready, skipped update")
		}
	}
}

func notifyNewPacket() {
	log.Debug("New packet received, triggering client update")
	updateClients()
}


func pixelsToJSON(pixels [][]bool) string {
	jsonPixels := make([][]int, len(pixels))
	for i, row := range pixels {
		jsonPixels[i] = make([]int, len(row))
		for j, pixel := range row {
			if pixel {
				jsonPixels[i][j] = 1
			} else {
				jsonPixels[i][j] = 0
			}
		}
	}
	jsonBytes, err := json.Marshal(jsonPixels)
	if err != nil {
		log.Errorf("Error marshaling pixels to JSON: %v", err)
		return "[]"
	}
	return string(jsonBytes)
}
