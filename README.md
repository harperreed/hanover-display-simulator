# Hanover Display Simulator

## Overview

The Hanover Display Simulator is a Go application that simulates a Hanover flipdot display. It listens for packets on a serial port, processes the data, and visualizes the display state through a web interface. This simulator is useful for testing and developing applications that interact with Hanover displays without needing the physical hardware.

## Features

- Simulates a Hanover flipdot display with configurable dimensions
- Listens for data on a serial port
- Processes incoming packets according to the Hanover protocol
- Provides a web interface to visualize the current display state
- Includes a test simulator for generating sample data
- Offers endpoints for retrieving packet history and raw display data

## Prerequisites

- Go 1.15 or higher
- Git
- socat (for creating virtual serial ports)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/hanover-display-simulator.git
   cd hanover-display-simulator
   ```

2. Install the required dependencies:
   ```
   go get github.com/gin-gonic/gin
   go get github.com/sirupsen/logrus
   go get github.com/tarm/serial
   go get gopkg.in/yaml.v2
   ```

3. Install socat (if not already installed):
   - On Ubuntu/Debian: `sudo apt-get install socat`
   - On macOS with Homebrew: `brew install socat`
   - On other systems, refer to the socat documentation for installation instructions.

## Configuration

Create a `config.yaml` file in the project root directory with the following content:

```yaml
columns: 96
rows: 16
address: 1
serial_port: "/dev/pts/2"  # Update this to match your system
baud_rate: 4800
web_port: ":8080"
```

Adjust the values as needed to match your desired configuration.

## Setting up Virtual Serial Ports

To test the simulator without physical hardware, you can use socat to create a pair of virtual serial ports:

1. Open a terminal and run the following command:
   ```
   socat -d -d pty,raw,echo=0 pty,raw,echo=0
   ```

2. socat will output two device names, for example:
   ```
   2023/04/01 12:00:00 socat[1234] N PTY is /dev/pts/2
   2023/04/01 12:00:00 socat[1234] N PTY is /dev/pts/3
   ```

3. Update the `serial_port` in your `config.yaml` to use one of these devices (e.g., `/dev/pts/2`).

4. You can use the other device (e.g., `/dev/pts/3`) to send test data to the simulator.

Keep this socat process running in the background while using the simulator.

## Usage

1. Start the simulator:
   ```
   go run main.go
   ```

2. The application will start and display log messages in the console.

3. Open a web browser and navigate to `http://localhost:8080` to see the visual representation of the display.

4. To send test data, the application includes a test simulator that sends packets periodically.

5. To stop the application, press Ctrl+C. The application will shut down gracefully.

## Web Endpoints

- `/`: Displays a visual representation of the current display state
- `/packets`: Returns a JSON array of received packet metadata
- `/display`: Returns a JSON representation of the current display state

## Sending Custom Data

To send custom data to the simulator, you can use any serial communication tool to send packets to the configured serial port. For example, using the second virtual serial port created by socat:

```
echo -ne "\x02\x11\x01\x00\xC0\xAA\xAA\xAA...\x03" > /dev/pts/3
```

The packet format should follow the Hanover protocol:

```
[STX (0x02)][Command (0x1)][Address (0x1)][Resolution (2 bytes)][Pixel Data][ETX (0x03)][Checksum (2 bytes)]
```

## Development

The main components of the application are:

- `main.go`: Contains the main application logic
- `config.yaml`: Configuration file for the application

To modify the behavior of the simulator, you can edit the `parseData` function in `main.go`.

## Troubleshooting

If you encounter issues:

1. Check that the serial port in `config.yaml` matches the device created by socat.
2. Ensure that socat is running and the virtual serial ports are created successfully.
3. Verify that no other application is using the specified serial port.
4. Check that the incoming data matches the expected Hanover protocol format.
5. If the simulator is not receiving data, try sending data to the other virtual serial port created by socat.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
