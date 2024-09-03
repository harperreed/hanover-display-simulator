package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestReassemblePacket(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected [][]byte
	}{
		{
			name:     "Complete packet",
			input:    []byte{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0x03, 0x00, 0x00},
			expected: [][]byte{{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0x03, 0x00, 0x00}},
		},
		{
			name:     "Partial packet",
			input:    []byte{0x02, 0x11, 0x01, 0x00, 0xC0},
			expected: [][]byte{},
		},
		{
			name:     "Multiple complete packets",
			input:    []byte{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0x03, 0x00, 0x00, 0x02, 0x11, 0x01, 0x00, 0xC0, 0xBB, 0x03, 0x00, 0x00},
			expected: [][]byte{{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0x03, 0x00, 0x00}, {0x02, 0x11, 0x01, 0x00, 0xC0, 0xBB, 0x03, 0x00, 0x00}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := reassemblePacket(tc.input)
			if !bytesSliceEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestParseData(t *testing.T) {
	// Initialize the display for testing
	config = Config{
		Columns: 96,
		Rows:    16,
		Address: 1,
	}
	initializeDisplay()

	testCases := []struct {
		name           string
		input          []byte
		expectedPixels int
	}{
		{
			name:           "Valid packet",
			input:          []byte{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0xAA, 0x03, 0x00, 0x00},
			expectedPixels: 8, // 0xAA in binary is 10101010, so 4 pixels per byte
		},
		{
			name:           "Invalid start byte",
			input:          []byte{0x03, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0xAA, 0x03, 0x00, 0x00},
			expectedPixels: 0,
		},
		{
			name:           "Invalid end byte",
			input:          []byte{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0xAA, 0x02, 0x00, 0x00},
			expectedPixels: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parseData(tc.input)
			updatedPixels := countUpdatedPixels()
			if updatedPixels != tc.expectedPixels {
				t.Errorf("Expected %d updated pixels, got %d", tc.expectedPixels, updatedPixels)
			}
		})
	}
}

func TestLogPacketToFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "packet_log_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Set the logFile to our temporary file
	logFile = tmpfile

	testPacket := Packet{
		Timestamp: time.Now(),
		Data:      []byte{0x02, 0x11, 0x01, 0x00, 0xC0, 0xAA, 0x03, 0x00, 0x00},
	}

	err = logPacketToFile(testPacket)
	if err != nil {
		t.Fatalf("Failed to log packet to file: %v", err)
	}

	// Read the contents of the file
	fileContents, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Unmarshal the JSON data
	var loggedPacket Packet
	err = json.Unmarshal(fileContents, &loggedPacket)
	if err != nil {
		t.Fatalf("Failed to unmarshal logged packet: %v", err)
	}

	// Compare the logged packet with the original
	if !bytes.Equal(loggedPacket.Data, testPacket.Data) {
		t.Errorf("Logged packet data does not match original. Expected %v, got %v", testPacket.Data, loggedPacket.Data)
	}
}

func bytesSliceEqual(a, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func countUpdatedPixels() int {
	count := 0
	for _, row := range display.pixels {
		for _, pixel := range row {
			if pixel {
				count++
			}
		}
	}
	return count
}
