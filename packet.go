package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

type Packet struct {
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}

var (
	packetLog     []Packet
	packetChan    = make(chan Packet, 100)
	partialPacket []byte
	logFile       *os.File
	logMutex      sync.Mutex
)

func initPacketLogging() error {
	var err error
	logFile, err = os.OpenFile("packet_log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return nil
}

func closePacketLogging() {
	if logFile != nil {
		logFile.Close()
	}
}

func logPacketToFile(packet Packet) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	jsonPacket, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	if _, err := logFile.Write(append(jsonPacket, '\n')); err != nil {
		return err
	}

	return nil
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

		// Log packet to JSON file
		if err := logPacketToFile(packet); err != nil {
			log.Errorf("Failed to log packet to file: %v", err)
		}

		parseData(packet.Data)
		notifyNewPacket() // Notify clients about the new packet
	}
}

func parseData(data []byte) {
    log.Infof("Parsing data: length=%d", len(data))
    if len(data) < 9 { // Minimum length: STX + CMD + ADDR + RES(2) + ETX + CHECKSUM(2)
        log.Warn("Received data too short")
        return
    }

    if data[0] != 0x02 || data[len(data)-3] != 0x03 {
        log.Warn("Invalid start or end byte")
        return
    }

    // Parse command (we're not using it currently, but it might be useful later)
    command := data[1]
    log.Infof("Command: %X", command)

    // Parse address
    addressStr := string(data[2])
    address, err := strconv.Atoi(addressStr)
    if err != nil {
        log.Warnf("Error parsing address: %v", err)
        return
    }
    if address != config.Address {
        log.Warnf("Message not for this display. Expected: %d, Got: %d", config.Address, address)
        return
    }

    // Parse resolution
    resolutionStr := string(data[3:5])
    resolution, err := strconv.ParseUint(resolutionStr, 16, 16)
    if err != nil {
        log.Warnf("Error parsing resolution: %v", err)
        return
    }
    expectedResolution := uint64((config.Rows * config.Columns) / 8)
    if resolution != expectedResolution {
        log.Warnf("Unexpected resolution. Expected: %d, Got: %d", expectedResolution, resolution)
    }

    // Parse pixel data
    pixelData := data[5 : len(data)-3]
    log.Infof("Pixel data length: %d", len(pixelData))

    updatedPixels := updateDisplay(pixelData)

    log.Infof("Data parsed successfully. Updated %d pixels.", updatedPixels)

    // Log the first few rows of the display for debugging
    for i := 0; i < min(5, len(display.pixels)); i++ {
        log.Infof("Row %d: %v", i, display.pixels[i][:min(10, len(display.pixels[i]))])
    }
}

func reassemblePacket(data []byte) [][]byte {
    var completePackets [][]byte

    // Append new data to any existing partial packet
    partialPacket = append(partialPacket, data...)

    for len(partialPacket) > 0 {
        // Find start byte
        startIndex := bytes.IndexByte(partialPacket, 0x02)
        if startIndex == -1 {
            // No start byte found, clear partial packet
            partialPacket = nil
            return completePackets
        }

        // Remove any data before the start byte
        partialPacket = partialPacket[startIndex:]

        // Find end byte
        endIndex := bytes.IndexByte(partialPacket, 0x03)
        if endIndex == -1 || len(partialPacket) < endIndex+3 {
            // End byte not found or not enough data for checksum, keep accumulating
            return completePackets
        }

        // We have a complete packet
        completePacket := partialPacket[:endIndex+3]
        completePackets = append(completePackets, completePacket)

        // Remove the complete packet from partial data
        partialPacket = partialPacket[endIndex+3:]
    }

    return completePackets
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
