package main

import (
    "testing"
    "time"
)

func TestTestSimulator(t *testing.T) {
    // Set up a test configuration
    config = Config{
        Columns: 96,
        Rows:    16,
        Address: 1,
    }

    // Initialize the display and packet channel
    initializeDisplay()
    packetChan = make(chan Packet, 100)

    // Start a goroutine to process packets
    go func() {
        for packet := range packetChan {
            parseData(packet.Data)
        }
    }()

    // Run the test simulator
    go testSimulator()

    // Wait for a short time to allow the simulator to run and packets to be processed
    time.Sleep(200 * time.Millisecond)

    // Check if a packet was sent
    if len(packetChan) == 0 {
        t.Error("No packets were sent by the simulator")
    }

    // Check if any packets were processed
    if len(packetLog) == 0 {
        t.Error("No packets were processed")
    }

    // Log the first few rows of the display for debugging
    display.mu.Lock()
    defer display.mu.Unlock()
    t.Log("Current display state:")
    for i := 0; i < min(5, len(display.pixels)); i++ {
        t.Logf("Row %d: %v", i, display.pixels[i][:min(10, len(display.pixels[i]))])
    }

    // Count the number of set pixels
    pixelsSet := 0
    for _, row := range display.pixels {
        for _, pixel := range row {
            if pixel {
                pixelsSet++
            }
        }
    }
    t.Logf("Total pixels set: %d", pixelsSet)

    // Check if the expected number of pixels are set
    if pixelsSet != 192 {
        t.Errorf("Expected 192 pixels to be set, got %d", pixelsSet)
    }
}
