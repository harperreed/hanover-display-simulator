package main

import (
	"strconv"
	"sync"
)

type HanoverDisplay struct {
	pixels [][]bool
	mu     sync.Mutex
}

var display HanoverDisplay

func initializeDisplay() {
	display = HanoverDisplay{
		pixels: make([][]bool, config.Rows),
	}
	for i := range display.pixels {
		display.pixels[i] = make([]bool, config.Columns)
	}
}

func updateDisplay(pixelData []byte) int {
    display.mu.Lock()
    defer display.mu.Unlock()

    updatedPixels := 0
    byteIndex := 0
    for col := 0; col < config.Columns; col++ {
        if byteIndex+1 >= len(pixelData) {
            break
        }
        byteVal, err := strconv.ParseUint(string(pixelData[byteIndex:byteIndex+2]), 16, 8)
        if err != nil {
            log.Warnf("Error parsing pixel data at position %d: %v", byteIndex, err)
            continue
        }
        for bit := 0; bit < 8; bit++ {
            row := bit
            if row < config.Rows {
                newValue := (byte(byteVal) & (1 << uint(7-bit))) != 0
                if display.pixels[row][col] != newValue {
                    display.pixels[row][col] = newValue
                    updatedPixels++
                }
            }
        }
        byteIndex += 2
    }

    return updatedPixels
}
