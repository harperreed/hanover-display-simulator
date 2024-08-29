package main

import (
	"html/template"
	"net/http"
	"time"

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
</head>
<body>
    <h1>Hanover Display Simulator</h1>
    <table>
        {{range .}}
        <tr>
            {{range .}}
            <td class="{{if .}}on{{else}}off{{end}}"></td>
            {{end}}
        </tr>
        {{end}}
    </table>
</body>
</html>
`

func runWebServer() {
	r := gin.Default()
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
	r.GET("/", func(c *gin.Context) {
		display.mu.Lock()
		defer display.mu.Unlock()

		tmpl, err := template.New("display").Parse(htmlTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error parsing template")
			return
		}

		err = tmpl.Execute(c.Writer, display.pixels)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error executing template")
			return
		}
	})
	if err := r.Run(config.WebPort); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
