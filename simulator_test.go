package main

import (
	"testing"
	"time"
)

func TestTestSimulator(t *testing.T) {
    config = Config{
        Columns: 96,
        Rows:    16,
        Address: 1,
    }

    initializeDisplay()
    packetChan = make(chan Packet, 100)

    go testSimulator()

    time.Sleep(100 * time.Millisecond)

    select {
    case packet := <-packetChan:
        if packet.Data[0] != 0x02 || packet.Data[len(packet.Data)-1] != 0x03 {
            t.Error("Invalid packet structure")
        }
        if len(packet.Data) != 198 {
            t.Errorf("Unexpected packet length: got %d, want 198", len(packet.Data))
        }
    default:
        t.Error("No packet received from simulator")
    }

    display.mu.Lock()
    defer display.mu.Unlock()

    pixelsSet := 0
    for _, row := range display.pixels {
        for _, pixel := range row {
            if pixel {
                pixelsSet++
            }
        }
    }

    if pixelsSet != 192 {
        t.Errorf("Expected 192 pixels to be set, got %d", pixelsSet)
    }
}
