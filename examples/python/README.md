# Project Overview

This repository contains a set of Python scripts designed to interact with flip-dot displays using a serial communication protocol. The primary functionality revolves around creating images, sending commands, and logging operations to a file for further analysis and debugging.

## Directory Structure

The repository is organized as follows:

```
python/
├── checkerboard.py            # Script to create a checkerboard image and send it to the display
├── proxy.log                  # Log file to record communication commands and responses
├── serial_logger_tester.py    # Sends example packets to the display for testing
├── serial_proxy.py            # Provides a wrapper for serial communication with logging functionality
├── test.py                    # Initializes the display and starts tests
└── test_pattern.py            # Similar to test.py, used for testing patterns on the display
```

## File Descriptions

### 1. `checkerboard.py`
This script generates a checkerboard pattern image and sends it to the connected flip-dot display. It uses the `HanoverController` and `HanoverSign` classes from the `pyflipdot` package to manage display operations effectively.

### 2. `proxy.log`
This is a log file that captures all write operations sent to the flip-dot display. Each entry includes the timestamp, operation type, and the data sent. This log helps to trace the communication history and diagnose any operational issues.

### 3. `serial_logger_tester.py`
The `serial_logger_tester.py` script is used to repeatedly send a predefined test packet to the display every second. This can be useful for validating the communication setup and testing how the display responds to incoming data.

### 4. `serial_proxy.py`
This file defines the `SerialProxy` class, which encapsulates serial communication functionality. It handles reading and writing data over serial ports while logging operations to a specified file. This class is essential for maintaining clean and structured logging of all serial interactions.

### 5. `test.py`
The `test.py` script sets up the connection to the flip-dot display and starts a test running predefined signs. It's a straightforward implementation aimed at ensuring that the display can receive and render sign data correctly.

### 6. `test_pattern.py`
Similar to `test.py`, this script also initializes the display and runs tests, but it may include different patterns or tests as needed for validation purposes.

## Requirements

- Python 3.x
- Required libraries:
  - `numpy`
  - `pyflipdot`
  - `pyserial`
  
## Usage

To use this project, ensure that your environment is set up correctly with the necessary dependencies installed. You can then run any of the provided scripts to interact with the flip-dot display. Make sure to specify the correct serial port in the scripts.

## Logging

All operations sent to the display are logged in `proxy.log`. This log captures detailed information about commands issued to the display, which can be used for debugging and performance monitoring.

## Conclusion

This repository serves as a foundational tool for managing flip-dot displays via serial communication. With logging capabilities and examples for sending various data formats, it can be adapted or expanded for a variety of applications involving flip-dot technology.
