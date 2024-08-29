# README for the Main Directory

## Overview

This repository is designed for a Python application focused on controlling Hanover LED displays through serial communication. The application allows users to create custom images, send test patterns, and log all communication activities. This README provides an overview of the contents available in this directory.

## Directory Structure

Within the main directory, the following files are available:

```
python/
├── README.md                    # Project overview and instructions.
├── checkerboard.py              # Script to create and send a checkerboard image to the display.
├── output.txt                   # Documentation file outlining the repository's structure and file contents.
├── proxy.log                     # Log file for recording all serial communication activities.
├── serial_logger_tester.py      # Script for testing serial communication via a sample packet.
├── serial_proxy.py              # Implementation of a serial communication proxy with logging capabilities.
├── test.py                      # Script to initialize and test the Hanover sign controller.
└── test_pattern.py              # Script that initializes the sign controller and runs test patterns.
```

## File Descriptions

- **README.md**: This file provides an overview of the project, including the purpose, installation instructions, usage guidelines, logging formats, and contribution information.

- **checkerboard.py**: A Python script that uses the `SerialProxy` and `HanoverController` classes to generate a checkerboard image and send it to the Hanover LED display.

- **output.txt**: A text file containing detailed documentation of the repository's structure and file contents. This file serves as a reference for understanding the project organization.

- **proxy.log**: A log file that captures all serial communication activities with the LED display. The log entries include timestamps, the type of operations (like READ or WRITE), and the data transmitted or received in hexadecimal format.

- **serial_logger_tester.py**: A utility script designed for testing the serial communication with the window sign by continuously sending a predefined sample packet at set intervals.

- **serial_proxy.py**: This file implements the `SerialProxy` class, which provides methods for managing serial communication and logging, encapsulating basic operations like reading, writing, flushing, and logging data transactions in a structured format.

- **test.py**: A straightforward script that initializes the Hanover sign controller and executes test signs for ensuring the display works as expected.

- **test_pattern.py**: Similar to `test.py`, this script initializes the sign controller and runs predefined test patterns on the display for validation purposes.

## Installation and Usage

Refer to the `README.md` for detailed installation instructions and usage guidelines on how to run the individual scripts to interact with the Hanover LED display. 

This repository is designed for developers and engineers interested in working with LED display systems and serial communication protocols. Feel free to explore, utilize the scripts, and contribute as needed!
