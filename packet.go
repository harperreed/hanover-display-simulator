package main

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
	"time"
 	// "encoding/hex"
    "strconv"
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
    if len(data) < 9 {
        log.Warnf("Received data too short: got %d bytes, expected at least 9", len(data))
        return
    }

    log.Infof("Packet header: %X", data[:5])
    log.Infof("Packet footer: %X", data[len(data)-3:])
    log.Infof("Full packet: %X", data)

    if data[0] != 0x02 || data[len(data)-3] != 0x03 {
        log.Warn("Invalid start or end byte")
        return
    }

    if !verifyChecksum(data) {
        log.Warn("Invalid checksum")
        return
    }

    // Parse command
    command := data[1]
    log.Infof("Command: %X", command)

    // Parse address
    address := int(data[2])
    if address != config.Address {
        log.Warnf("Message not for this display. Expected: %d, Got: %d", config.Address, address)
        return
    }

    // Parse resolution
    resolution := uint16(data[3])<<8 | uint16(data[4])
    expectedResolution := uint16((config.Rows * config.Columns) / 8)
    if resolution != expectedResolution {
        log.Warnf("Unexpected resolution. Expected: %d, Got: %d", expectedResolution, resolution)
    }

    // Parse pixel data
    pixelData := data[5 : len(data)-3]
    log.Infof("Pixel data length: %d", len(pixelData))
    log.Infof("First few bytes of pixel data: %X", pixelData[:min(20, len(pixelData))])

    updatedPixels := updateDisplay(pixelData)

    log.Infof("Data parsed successfully. Updated %d pixels.", updatedPixels)

    // Log the first few rows of the display for debugging
    for i := 0; i < min(5, len(display.pixels)); i++ {
        log.Infof("Row %d: %v", i, display.pixels[i][:min(20, len(display.pixels[i]))])
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

func verifyChecksum(data []byte) bool {
    if len(data) < 4 {
        log.Warn("Packet too short for checksum verification")
        return false
    }

    log.Infof("Full packet data: %X", data)

    // Extract the checksum from the packet
    checksumStr := string(data[len(data)-2:])
    log.Infof("Checksum string: %s", checksumStr)

    // Parse the checksum as a hexadecimal number
    packetChecksum, err := strconv.ParseUint(checksumStr, 16, 8)
    if err != nil {
        log.Errorf("Failed to parse checksum as hex: %v", err)
        return false
    }
    log.Infof("Packet checksum: %02X", packetChecksum)

    // Calculate the sum of all bytes excluding SOT (first byte) and the checksum itself
    var sum uint32
    for _, b := range data[1 : len(data)-2] {
        sum += uint32(b)
    }
    log.Infof("Sum of data (excluding SOT and checksum): %08X", sum)

    // Subtract SOT (0x02) from the sum
    sum -= 0x02
    log.Infof("Sum after subtracting SOT: %08X", sum)

    // XOR with 0xFF (255)
    sum ^= 0xFF
    log.Infof("Sum after XOR with 0xFF: %08X", sum)

    // Add 1
    sum += 1
    log.Infof("Sum after adding 1: %08X", sum)

    // Take the least significant byte
    calculatedChecksum := byte(sum & 0xFF)
    log.Infof("Calculated checksum: %02X", calculatedChecksum)

    // Compare the calculated checksum with the packet checksum
    isValid := calculatedChecksum == byte(packetChecksum)
    log.Infof("Checksum valid: %v", isValid)

    return isValid
}

// CRC-16 (CCITT) implementation
func crc16(data []byte) uint16 {
    crc := uint16(0xFFFF)
    for _, b := range data {
        crc ^= uint16(b) << 8
        for i := 0; i < 8; i++ {
            if crc&0x8000 != 0 {
                crc = (crc << 1) ^ 0x1021
            } else {
                crc <<= 1
            }
        }
    }
    return crc
}
