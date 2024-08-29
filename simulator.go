package main

import (
	"time"
)

func testSimulator() {
	log.Info("Running test simulation")
	testPacket := []byte{0x02, 0x11, 0x01, 0x00, 0xC0}
	for i := 0; i < 192; i++ {
		testPacket = append(testPacket, 0xAA)
	}
	testPacket = append(testPacket, 0x03)

	packet := Packet{
		Timestamp: time.Now(),
		Data:      testPacket,
	}
	packetChan <- packet
	log.Infof("Sent test packet: length=%d", len(testPacket))

	time.Sleep(time.Second) // Give some time for processing

	display.mu.Lock()
	defer display.mu.Unlock()
	log.Info("Current display state:")
	for i := 0; i < min(5, len(display.pixels)); i++ {
		log.Infof("Row %d: %v", i, display.pixels[i][:min(10, len(display.pixels[i]))])
	}
}
