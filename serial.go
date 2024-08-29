package main

import (
	"time"

	"github.com/tarm/serial"
)

func readSerialPort() {
	serialConfig := &serial.Config{
		Name: config.SerialPort,
		Baud: config.BaudRate,
	}
	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatalf("Error opening serial port: %v", err)
	}

	log.Infof("Started reading from serial port %s", config.SerialPort)

	for {
		buf := make([]byte, 512)
		n, err := port.Read(buf)
		if err != nil {
			log.Errorf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			data := buf[:n]
			log.Infof("Received data: length=%d, first byte=0x%02X, last byte=0x%02X",
				len(data), data[0], data[len(data)-1])

			completePackets := reassemblePacket(data)
			for _, completePacket := range completePackets {
				packet := Packet{
					Timestamp: time.Now(),
					Data:      completePacket,
				}
				packetChan <- packet
				log.Infof("Assembled complete packet: length=%d, first byte=0x%02X, last byte=0x%02X",
					len(completePacket), completePacket[0], completePacket[len(completePacket)-1])
			}
		}
	}
}