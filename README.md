```markdown
# 🖥️ Hanover Display Simulator

## 📋 Overview

Welcome to the **Hanover Display Simulator**! This Go application is designed to emulate a Hanover flipdot display. It connects to a serial port to receive packets, processes this data, and visualizes the display's state via a web interface. This tool is essential for testing and developing applications that interact with actual Hanover displays without having the physical hardware at hand.

## 🎉 Features

- 🎨 Simulates a Hanover flipdot display with customizable dimensions.
- 📡 Listens for data over a specified serial port.
- 📜 Processes incoming packets following the Hanover display protocol.
- 🌐 Provides a web interface to visualize the live display state.
- 🧪 Incorporates a test simulator for generating sample display data.
- 🔍 Offers endpoints to retrieve packet history and raw display data.

## 👨‍💻 How to Use

### 1. Prerequisites

Before running the simulator, ensure you have the following requirements:

- Go 1.15 or higher
- Git
- `socat` for creating virtual serial ports

### 2. Installation

Follow these steps to clone and set up the project:

```bash
# Clone the repository
git clone https://github.com/harperreed/hanover-display-simulator.git
cd hanover-display-simulator

# Install dependencies
go get github.com/gin-gonic/gin
go get github.com/sirupsen/logrus
go get github.com/tarm/serial
go get gopkg.in/yaml.v2

# Install socat (if you haven't already)
# For Ubuntu/Debian:
sudo apt-get install socat
# For macOS with Homebrew:
brew install socat
```

### 3. Configure the Simulator

Create a `config.yaml` file in the project root with the following content:

```yaml
columns: 96
rows: 16
address: 1
serial_port: "/dev/pts/2"  # Update this to match your system
baud_rate: 4800
web_port: ":8080"
```

Adjust the values based on your system's setup.

### 4. Setting up Virtual Serial Ports

To test the simulator without actual hardware, use `socat` to create virtual serial ports:

1. Run the following command in your terminal:

    ```bash
    socat -d -d pty,raw,echo=0 pty,raw,echo=0
    ```

2. Update the `serial_port` in your `config.yaml` with one of the device names output by `socat`.

3. Keep the `socat` command running while you use the simulator.

### 5. Start the Simulator

To launch the simulator, execute:

```bash
go run main.go
```

### 6. View the Display

Open a web browser and navigate to `http://localhost:8080` to view the display's visual representation. To send test data, use the second virtual serial port created by `socat`.

## 🚀 Tech Info

This project is built using the Go programming language. Some of the key components include:

- **Main Logic:** The core logic resides in `main.go`, which handles the application flow.
- **Serial Communication:** Implemented via the `serial.go`, enabling the reading of incoming packets.
- **Display Management:** Managed by `display.go`, this file maintains and updates the display state.
- **Web Server:** The `webserver.go` file provides a web interface using the Gin web framework, serving display updates and packet information.
- **Configuration Management:** Configuration is handled via `config.yaml` and `config.go`, ensuring easy adjustments.

### Dependencies

- `gin-gonic/gin` for web server capabilities
- `sirupsen/logrus` for structured logging
- `tarm/serial` for serial port communication
- `gopkg.in/yaml.v2` for YAML file parsing

## 🙌 Contributing

Contributions are welcome! If you have suggestions for improvements or features, feel free to submit a pull request!

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Happy Coding! 🌍
```

