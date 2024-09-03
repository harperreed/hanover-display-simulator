package main

import (
	"reflect"
	"sync"
	"fmt"
	"testing"
)

func TestUpdateDisplay(t *testing.T) {
    tests := []struct {
        name             string
        initialDisplay   [][]bool
        pixelData        []byte
        expectedDisplay  [][]bool
        expectedUpdates  int
    }{
        {
            name:           "Empty display, no changes",
            initialDisplay: makeEmptyDisplay(),
            pixelData:      []byte("0000000000000000000000000000000000000000000000000000000000000000000000"),
            expectedDisplay: makeEmptyDisplay(),
            expectedUpdates: 0,
        },
        {
            name:           "Set all pixels to on",
            initialDisplay: makeEmptyDisplay(),
            pixelData:      []byte("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
            expectedDisplay: makeFullDisplay(true),
            expectedUpdates: 16 * 56,
        },
        {
            name:           "Set all pixels to off",
            initialDisplay: makeFullDisplay(true),
            pixelData:      []byte("0000000000000000000000000000000000000000000000000000000000000000000000"),
            expectedDisplay: makeEmptyDisplay(),
            expectedUpdates: 16 * 56,
        },
        {
            name:           "Set alternate columns",
            initialDisplay: makeEmptyDisplay(),
            pixelData:      []byte("FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00"),
            expectedDisplay: makeAlternateColumnDisplay(),
            expectedUpdates: 16 * 28,
        },
        {
            name:           "Invalid data",
            initialDisplay: makeEmptyDisplay(),
            pixelData:      []byte("INVALID"),
            expectedDisplay: makeEmptyDisplay(),
            expectedUpdates: 0,
        },
        {
            name:           "Partial update",
            initialDisplay: makeEmptyDisplay(),
            pixelData:      []byte("FFFF"),
            expectedDisplay: makePartiallyUpdatedDisplay(),
            expectedUpdates: 16,
        },
    }

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            display = HanoverDisplay{
                pixels: tt.initialDisplay,
                mu:     sync.Mutex{},
            }

            fmt.Printf("Running test: %s\n", tt.name)
            fmt.Printf("Input pixel data: %X\n", tt.pixelData)

            updatedPixels := updateDisplay(tt.pixelData)

            if updatedPixels != tt.expectedUpdates {
                t.Errorf("Expected %d updated pixels, got %d", tt.expectedUpdates, updatedPixels)
            }

            if !reflect.DeepEqual(display.pixels, tt.expectedDisplay) {
                t.Errorf("Display state doesn't match expected state")
                t.Logf("Expected:\n%v", tt.expectedDisplay)
                t.Logf("Got:\n%v", display.pixels)
            }

            fmt.Printf("Test completed: %s\n\n", tt.name)
        })
    }
}
func makeEmptyDisplay() [][]bool {
	display := make([][]bool, config.Rows)
	for i := range display {
		display[i] = make([]bool, config.Columns)
	}
	return display
}

func makeFullDisplay(value bool) [][]bool {
	display := make([][]bool, config.Rows)
	for i := range display {
		display[i] = make([]bool, config.Columns)
		for j := range display[i] {
			display[i][j] = value
		}
	}
	return display
}

func makeAlternateColumnDisplay() [][]bool {
	display := makeEmptyDisplay()
	for col := 0; col < config.Columns; col += 2 {
		for row := 0; row < config.Rows; row++ {
			display[row][col] = true
		}
	}
	return display
}

func makePartiallyUpdatedDisplay() [][]bool {
	display := makeEmptyDisplay()
	for row := 0; row < config.Rows; row++ {
		display[row][0] = true
	}
	return display
}
