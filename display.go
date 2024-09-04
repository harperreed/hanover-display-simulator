package main

import (
	"strconv"
	"sync"
	"fmt"
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

// original code
// func updateDisplay(pixelData []byte) int {
//     display.mu.Lock()
//     defer display.mu.Unlock()

//     updatedPixels := 0
//     byteIndex := 0

//     for col := 0; col < config.Columns; col++ {
//         for rowByte := 0; rowByte < (config.Rows + 7) / 8; rowByte++ {
//             if byteIndex+1 >= len(pixelData) {
//                 return updatedPixels
//             }
//             byteVal, err := strconv.ParseUint(string(pixelData[byteIndex:byteIndex+2]), 16, 8)
//             if err != nil {
//                 log.Warnf("Error parsing pixel data at position %d: %v", byteIndex, err)
//                 byteIndex += 2
//                 continue
//             }
//             for bit := 0; bit < 8; bit++ {
//                 row := rowByte*8 + bit
//                 if row < config.Rows {
//                     newValue := (byte(byteVal) & (1 << uint(7-bit))) != 0
//                     if display.pixels[row][col] != newValue {
//                         display.pixels[row][col] = newValue
//                         updatedPixels++
//                     }
//                 }
//             }
//             byteIndex += 2
//         }
//     }

//     return updatedPixels
// }

//try two
// func updateDisplay(pixelData []byte) int {
//     display.mu.Lock()
//     defer display.mu.Unlock()
//     updatedPixels := 0
//     byteIndex := 0
//     for col := 0; col < config.Columns; col++ {
//         for rowByte := 0; rowByte < (config.Rows + 7) / 8; rowByte++ {
//             if byteIndex+1 >= len(pixelData) {
//                 return updatedPixels
//             }
//             byteVal, err := strconv.ParseUint(string(pixelData[byteIndex:byteIndex+2]), 16, 8)
//             if err != nil {
//                 log.Printf("Error parsing pixel data at position %d: %v", byteIndex, err)
//                 byteIndex += 2
//                 continue
//             }
//             for bit := 0; bit < 8; bit++ {
//                 row := rowByte*8 + bit
//                 if row < config.Rows {
//                     newValue := (byte(byteVal) & (1 << uint(bit))) != 0
//                     if display.pixels[row][col] != newValue {
//                         display.pixels[row][col] = newValue
//                         updatedPixels++
//                     }
//                 }
//             }
//             byteIndex += 2
//         }
//     }
//     return updatedPixels
// }

//try three
// func updateDisplay(pixelData []byte) int {
//     display.mu.Lock()
//     defer display.mu.Unlock()
//     updatedPixels := 0
//     byteIndex := 0
//     for col := 0; col < config.Columns; col++ {
//         for rowByte := 0; rowByte < (config.Rows + 7) / 8; rowByte++ {
//             if byteIndex+1 >= len(pixelData) {
//                 return updatedPixels
//             }
//             byteVal, err := strconv.ParseUint(string(pixelData[byteIndex:byteIndex+2]), 16, 8)
//             if err != nil {
//                 log.Printf("Error parsing pixel data at position %d: %v", byteIndex, err)
//                 byteIndex += 2
//                 continue
//             }
//             for bit := 0; bit < 8; bit++ {
//                 row := rowByte*8 + bit
//                 if row < config.Rows {
//                     newValue := (byte(byteVal) & (1 << uint(7-bit))) != 0
//                     if display.pixels[row][col] != newValue {
//                         display.pixels[row][col] = newValue
//                         updatedPixels++
//                     }
//                 }
//             }
//             byteIndex += 2
//         }
//     }
//     return updatedPixels
// }

// try four
// func updateDisplay(pixelData []byte) int {
//     display.mu.Lock()
//     defer display.mu.Unlock()
//     updatedPixels := 0
//     byteIndex := 0

//     fmt.Printf("Updating display with pixel data of length: %d\n", len(pixelData))

//     for col := 0; col < config.Columns; col++ {
//         for rowByte := 0; rowByte < (config.Rows + 7) / 8; rowByte++ {
//             if byteIndex+1 >= len(pixelData) {
//                 fmt.Printf("Reached end of pixel data at byteIndex: %d\n", byteIndex)
//                 return updatedPixels
//             }
//             byteVal, err := strconv.ParseUint(string(pixelData[byteIndex:byteIndex+2]), 16, 8)
//             if err != nil {
//                 fmt.Printf("Error parsing pixel data at position %d: %v\n", byteIndex, err)
//                 byteIndex += 2
//                 continue
//             }
//             fmt.Printf("Processing column %d, rowByte %d, byteVal: %02X\n", col, rowByte, byteVal)
//             for bit := 0; bit < 8; bit++ {
//                 row := rowByte*8 + bit
//                 if row < config.Rows {
//                     newValue := (byte(byteVal) & (1 << uint(7-bit))) != 0
//                     if display.pixels[row][col] != newValue {
//                         display.pixels[row][col] = newValue
//                         updatedPixels++
//                         fmt.Printf("Updated pixel at row %d, col %d to %v\n", row, col, newValue)
//                     }
//                 }
//             }
//             byteIndex += 2
//         }
//     }
//     fmt.Printf("Total updated pixels: %d\n", updatedPixels)
//     return updatedPixels
// }


func updateDisplay(pixelData []byte) int {
    display.mu.Lock()
    defer display.mu.Unlock()
    updatedPixels := 0
    byteIndex := 0

    fmt.Printf("Initializing display update. Pixel data length: %d\n", len(pixelData))
    fmt.Printf("Display dimensions: Rows=%d, Columns=%d\n", config.Rows, config.Columns)

    for col := 0; col < config.Columns; col++ {
        for rowByte := 0; rowByte < (config.Rows+7)/8; rowByte++ {
            if byteIndex+1 >= len(pixelData) {
                fmt.Printf("Reached end of pixel data at byteIndex: %d\n", byteIndex)
                return updatedPixels
            }
            fmt.Printf("Processing byteIndex: %d\n", byteIndex)
            byteVal, err := strconv.ParseUint(string(pixelData[byteIndex:byteIndex+2]), 16, 8)
            if err != nil {
                fmt.Printf("Error parsing pixel data at byteIndex %d: %v\n", byteIndex, err)
                byteIndex += 2
                continue
            }
            fmt.Printf("Processing column %d, rowByte %d, byteVal: %02X\n", col, rowByte, byteVal)
            for bit := 0; bit < 8; bit++ {
                row := rowByte*8 + bit
                if row < config.Rows {
                    newValue := (byte(byteVal) & (1 << uint(7-bit))) != 0
                    fmt.Printf("Checking pixel at row %d, col %d. Old value: %v, New value: %v\n", row, col, display.pixels[row][col], newValue)
                    if display.pixels[row][col] != newValue {
                        display.pixels[row][col] = newValue
                        updatedPixels++
                        fmt.Printf("Updated pixel at row %d, col %d to %v\n", row, col, newValue)
                    }
                }
            }
            byteIndex += 2
        }
    }
    fmt.Printf("Display update complete. Total updated pixels: %d\n", updatedPixels)
    return updatedPixels
}
