package main

import (
	"bytes"
	"strconv"
	"time"
)

type Packet struct {
	Timestamp time.Time
	Data      []byte
}

var (
	packetLog     []Packet
	packetChan    = make(chan Packet, 100)
	partialPacket []byte
)

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
	command := string(data[1])
	log.Infof("Command: %s", command)

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

	// Parse resolution (we're not using it currently, but it might be useful later)
	resolution := string(data[3:5])
	log.Infof("Resolution: %s", resolution)

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

	for {
		// Find start byte
		startIndex := bytes.IndexByte(partialPacket, 0x02)
		if startIndex == -1 {
			// No start byte found, keep all data as partial
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

		// If no more data in partial packet, we're done
		if len(partialPacket) == 0 {
			return completePackets
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
