package main

import (
	"testing"
	"bytes"
)

func TestInitializeDisplay(t *testing.T) {
	config = Config{
		Columns: 96,
		Rows:    16,
	}

	initializeDisplay()

	if len(display.pixels) != config.Rows {
		t.Errorf("Expected %d rows, got %d", config.Rows, len(display.pixels))
	}

	for _, row := range display.pixels {
		if len(row) != config.Columns {
			t.Errorf("Expected %d columns, got %d", config.Columns, len(row))
		}
	}
}

func TestUpdateDisplay(t *testing.T) {
	config = Config{
		Columns: 96,
		Rows:    16,
	}

	initializeDisplay()

	testCases := []struct {
		name           string
		input          []byte
		expectedPixels int
	}{
		{
			name:           "All pixels on",
			input:          bytes.Repeat([]byte{0xFF}, 192),
			expectedPixels: 96 * 16,
		},
		{
			name:           "All pixels off",
			input:          bytes.Repeat([]byte{0x00}, 192),
			expectedPixels: 0,
		},
		{
			name:           "Alternating pixels",
			input:          bytes.Repeat([]byte{0xAA}, 192),
			expectedPixels: (96 * 16) / 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedPixels := updateDisplay(tc.input)
			if updatedPixels != tc.expectedPixels {
				t.Errorf("Expected %d updated pixels, got %d", tc.expectedPixels, updatedPixels)
			}

			// Verify the display state
			actualPixels := countUpdatedPixels()
			if actualPixels != tc.expectedPixels {
				t.Errorf("Display state incorrect. Expected %d pixels on, got %d", tc.expectedPixels, actualPixels)
			}
		})
	}
}
