package main

import (
	"fmt"
	"testing"
)

// Setup test environment
func TestMain(m *testing.M) {
	// Set up the correct configuration
	config = Config{
		Rows:    16, // Adjust to your actual display size
		Columns: 8,
	}

	fmt.Printf("Setting config: Rows=%d, Columns=%d\n", config.Rows, config.Columns)
	initializeDisplay() // Initialize the display with the correct config
	m.Run()             // Run tests
}

// Test edge case where pixelData is smaller than expected
func TestUpdateDisplayShortData(t *testing.T) {
	fmt.Println("Testing display update with short pixel data")
	// Test case: Short data
	shortPixelData := []byte("FF") // Single byte (FF) should turn all bits of the first row to true

	updatedPixels := updateDisplay(shortPixelData)

	if updatedPixels != 8 { // Expect 8 pixels to be updated for the first row
		t.Errorf("Expected 8 pixels to be updated with short data, but got %d", updatedPixels)
	}

	// Ensure only the first row is updated
	expectedPixels := [][]bool{
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, false},
	}

	for row := 0; row < len(expectedPixels); row++ {
		for col := 0; col < config.Columns; col++ {
			if display.pixels[row][col] != expectedPixels[row][col] {
				t.Errorf("Expected pixel at row %d, col %d to be %v, but got %v", row, col, expectedPixels[row][col], display.pixels[row][col])
			}
		}
	}

	// Verify no further rows are unexpectedly changed
	for row := len(expectedPixels); row < config.Rows; row++ {
		for col := 0; col < config.Columns; col++ {
			if display.pixels[row][col] {
				t.Errorf("Unexpectedly found a true pixel at row %d, col %d, but expected false", row, col)
			}
		}
	}
}
