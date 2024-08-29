package main

import (
	"html/template"
	"net/http"
	"time"
	"bytes"
    "io"
    "sync"
	"github.com/gin-gonic/gin"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Hanover Display Simulator</title>
    <style>
        table {
            border-collapse: collapse;
        }
        td {
            width: 10px;
            height: 10px;
            border: 1px solid #444;
        }
        .on {
            background-color: yellow;
        }
        .off {
            background-color: black;
        }
    </style>
    <script>
        function setupEventSource() {
            var eventSource = new EventSource("/events");
            eventSource.onmessage = function(event) {
                document.getElementById("display-container").innerHTML = event.data;
            };
            eventSource.onerror = function(error) {
                console.error("EventSource failed:", error);
                eventSource.close();
                setTimeout(setupEventSource, 5000);
            };
        }
        window.onload = setupEventSource;
    </script>
</head>
<body>
    <h1>Hanover Display Simulator</h1>
    <div id="display-container">
        {{.}}
    </div>
</body>
</html>
`

const displayTableTemplate = `
<table>
    {{range .}}
    <tr>
        {{range .}}
        <td class="{{if .}}on{{else}}off{{end}}"></td>
        {{end}}
    </tr>
    {{end}}
</table>
`

var clients = make(map[chan string]bool)
var clientsMutex sync.Mutex

func runWebServer() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		display.mu.Lock()
		defer display.mu.Unlock()

		tmpl, err := template.New("display").Parse(htmlTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error parsing template")
			return
		}

		displayTmpl, err := template.New("displayTable").Parse(displayTableTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error parsing display table template")
			return
		}

		var displayHTML string
		buf := new(bytes.Buffer)
		err = displayTmpl.Execute(buf, display.pixels)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error executing template")
			return
		}
		displayHTML = buf.String()

		err = tmpl.Execute(c.Writer, displayHTML)
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
		c.JSON(http.StatusOK, display.pixels)
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

	displayTmpl, err := template.New("displayTable").Parse(displayTableTemplate)
	if err != nil {
		log.Errorf("Error parsing display table template: %v", err)
		return
	}

	var displayHTML string
	buf := new(bytes.Buffer)
	err = displayTmpl.Execute(buf, display.pixels)
	if err != nil {
		log.Errorf("Error executing template: %v", err)
		return
	}
	displayHTML = buf.String()

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
