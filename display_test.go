package main

import (
    "strings"
    "testing"
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

func TestUpdateDisplayWithVariousInputs(t *testing.T) {
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
        checkPattern   func([][]bool) bool
    }{
        {
            name:           "All pixels on",
            input:          []byte(strings.Repeat("FF", 24)),
            expectedPixels: 96 * 16,
            checkPattern: func(pixels [][]bool) bool {
                for _, row := range pixels {
                    for _, pixel := range row {
                        if !pixel {
                            return false
                        }
                    }
                }
                return true
            },
        },
        {
            name:           "All pixels off",
            input:          []byte(strings.Repeat("00", 24)),
            expectedPixels: 0,
            checkPattern: func(pixels [][]bool) bool {
                for _, row := range pixels {
                    for _, pixel := range row {
                        if pixel {
                            return false
                        }
                    }
                }
                return true
            },
        },
        {
            name:           "Alternating pixels",
            input:          []byte(strings.Repeat("AA", 24)),
            expectedPixels: (96 * 16) / 2,
            checkPattern: func(pixels [][]bool) bool {
                for i, row := range pixels {
                    for j, pixel := range row {
                        if pixel != ((i+j)%2 == 0) {
                            return false
                        }
                    }
                }
                return true
            },
        },
        {
            name:           "Alternating columns",
            input:          []byte(strings.Repeat("AA", 24)),
            expectedPixels: (96 * 16) / 2,
            checkPattern: func(pixels [][]bool) bool {
                for _, row := range pixels {
                    for j, pixel := range row {
                        if pixel != (j%2 == 0) {
                            return false
                        }
                    }
                }
                return true
            },
        },
        {
            name:           "Alternating rows",
            input:          []byte(strings.Repeat("FF00", 12)),
            expectedPixels: (96 * 16) / 2,
            checkPattern: func(pixels [][]bool) bool {
                for i, row := range pixels {
                    for _, pixel := range row {
                        if pixel != (i%2 == 0) {
                            return false
                        }
                    }
                }
                return true
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Reset display before each test
            initializeDisplay()

            updatedPixels := updateDisplay(tc.input)
            if updatedPixels != tc.expectedPixels {
                t.Errorf("Expected %d updated pixels, got %d", tc.expectedPixels, updatedPixels)
            }

            // Verify the display state
            actualPixels := countUpdatedPixelsInTest()
            if actualPixels != tc.expectedPixels {
                t.Errorf("Display state incorrect. Expected %d pixels on, got %d", tc.expectedPixels, actualPixels)
            }

            // Check the pattern
            if !tc.checkPattern(display.pixels) {
                t.Errorf("Display pattern is incorrect for test case: %s", tc.name)
            }

            // Print the first few rows for debugging
            for i := 0; i < minInt(5, len(display.pixels)); i++ {
                t.Logf("Row %d: %v", i, display.pixels[i][:minInt(10, len(display.pixels[i]))])
            }
        })
    }
}

// Helper function to count updated pixels
func countUpdatedPixelsInTest() int {
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

// Helper function to get the minimum of two integers
func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}
