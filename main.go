package main

import (
	// "encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Columns    int    `yaml:"columns"`
	Rows       int    `yaml:"rows"`
	Address    int    `yaml:"address"`
	SerialPort string `yaml:"serial_port"`
	BaudRate   int    `yaml:"baud_rate"`
	WebPort    string `yaml:"web_port"`
}

type HanoverDisplay struct {
	pixels [][]bool
	mu     sync.Mutex
}

type Packet struct {
	Timestamp time.Time
	Data      []byte
}

var (
	config     Config
	display    HanoverDisplay
	packetLog  []Packet
	packetChan = make(chan Packet, 100)
	log        = logrus.New()
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
	display = HanoverDisplay{
		pixels: make([][]bool, config.Rows),
	}
	for i := range display.pixels {
		display.pixels[i] = make([]bool, config.Columns)
	}

	go runWebServer()
	go processPackets()
	go readSerialPort()
	go testSimulator() // Run a test simulation

	// Keep the main goroutine running
	select {}
}

func loadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	return nil
}

func readSerialPort() {
	serialConfig := &serial.Config{
		Name: config.SerialPort,
		Baud: config.BaudRate,
	}
	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatalf("Error opening serial port: %v", err)
	}

	log.Infof("Started reading from serial port %s", config.SerialPort)

	for {
		buf := make([]byte, 256)
		n, err := port.Read(buf)
		if err != nil {
			log.Errorf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			data := buf[:n]
			packet := Packet{
				Timestamp: time.Now(),
				Data:      data,
			}
			packetChan <- packet
			log.Infof("Received packet: length=%d, first byte=0x%02X, last byte=0x%02X",
				len(data), data[0], data[len(data)-1])
		}
	}
}

func parseData(data []byte) {
	log.Infof("Parsing data: length=%d", len(data))
	if len(data) < 7 {
		log.Warn("Received data too short")
		return
	}

	if data[0] != 0x02 || data[len(data)-1] != 0x03 {
		log.Warn("Invalid start or end byte")
		return
	}

	addressByte := int(data[2])
	if addressByte != config.Address {
		log.Warnf("Message not for this display. Expected: %d, Got: %d", config.Address, addressByte)
		return
	}

	pixelData := data[5 : len(data)-1]
	log.Infof("Pixel data length: %d", len(pixelData))

	display.mu.Lock()
	defer display.mu.Unlock()

	updatedPixels := 0
	for i := 0; i < len(pixelData); i++ {
		byteVal := pixelData[i]
		col := i * 2 // Each byte represents 2 columns
		for bit := 0; bit < 8; bit++ {
			row := bit
			if col < config.Columns && row < config.Rows {
				newValue := (byteVal & (1 << uint(7-bit))) != 0
				if display.pixels[row][col] != newValue {
					display.pixels[row][col] = newValue
					updatedPixels++
				}
			}
			if col+1 < config.Columns && row < config.Rows {
				newValue := (byteVal & (1 << uint(3-bit))) != 0
				if display.pixels[row][col+1] != newValue {
					display.pixels[row][col+1] = newValue
					updatedPixels++
				}
			}
		}
	}
	log.Infof("Data parsed successfully. Updated %d pixels.", updatedPixels)
}

func processPackets() {
	log.Info("Started processing packets")
	for packet := range packetChan {
		packetLog = append(packetLog, packet)
		if len(packetLog) > 100 {
			packetLog = packetLog[1:]
		}
		log.Infof("Processing packet: timestamp=%v, length=%d",
			packet.Timestamp, len(packet.Data))
		parseData(packet.Data)
	}
}

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

func testSimulator() {
	log.Info("Running test simulation")
	testPacket := []byte{0x02, 0x11, 0x01, 0x00, 0xC0}
	for i := 0; i < 192; i++ {
		testPacket = append(testPacket, 0xAA)
	}
	testPacket = append(testPacket, 0x03)

	packet := Packet{
		Timestamp: time.Now(),
		Data:      testPacket,
	}
	packetChan <- packet
	log.Infof("Sent test packet: length=%d", len(testPacket))

	time.Sleep(time.Second) // Give some time for processing

	display.mu.Lock()
	defer display.mu.Unlock()
	log.Info("Current display state:")
	for i := 0; i < min(5, len(display.pixels)); i++ {
		log.Infof("Row %d: %v", i, display.pixels[i][:min(10, len(display.pixels[i]))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
