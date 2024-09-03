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

	// Run the test simulator
	go testSimulator()

	// Wait for a short time to allow the simulator to run
	time.Sleep(100 * time.Millisecond)

	// Check if a packet was sent
	select {
	case packet := <-packetChan:
		// Verify the packet structure
		if packet.Data[0] != 0x02 || packet.Data[len(packet.Data)-1] != 0x03 {
			t.Error("Invalid packet structure")
		}
		if len(packet.Data) != 197 { // STX + CMD + ADDR + RES(2) + 192 bytes of data + ETX
			t.Errorf("Unexpected packet length: got %d, want 197", len(packet.Data))
		}
	default:
		t.Error("No packet received from simulator")
	}

	// Verify the display state
	display.mu.Lock()
	defer display.mu.Unlock()

	// Check if any pixels are set (the exact pattern may vary, so we just check if some are set)
	pixelsSet := false
	for _, row := range display.pixels {
		for _, pixel := range row {
			if pixel {
				pixelsSet = true
				break
			}
		}
		if pixelsSet {
			break
		}
	}

	if !pixelsSet {
		t.Error("No pixels were set in the display")
	}
}
