import serial
import json
from datetime import datetime

class SerialProxy:
    def __init__(self, port, baudrate, log_file='serial_log.json'):
        self.serial = serial.Serial(port, baudrate)
        self.log_file = log_file

    def write(self, data):
        self.log("WRITE", data)
        return self.serial.write(data)

    def read(self, size=1):
        data = self.serial.read(size)
        self.log("READ", data)
        return data

    def readline(self):
        data = self.serial.readline()
        self.log("READLINE", data)
        return data

    def readlines(self):
        data = self.serial.readlines()
        self.log("READLINES", data)
        return data

    def in_waiting(self):
        return self.serial.in_waiting

    def flush(self):
        self.log("FLUSH", "")
        return self.serial.flush()

    def flushInput(self):
        self.log("FLUSH_INPUT", "")
        return self.serial.flushInput()

    def flushOutput(self):
        self.log("FLUSH_OUTPUT", "")
        return self.serial.flushOutput()

    def close(self):
        self.log("CLOSE", "")
        return self.serial.close()

    def open(self):
        self.log("OPEN", "")
        return self.serial.open()

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()

    @property
    def is_open(self):
        return self.serial.is_open

    def log(self, operation, data):
        log_entry = {
            "timestamp": datetime.now().isoformat(),
            "operation": operation,
            "data": data.hex() if isinstance(data, bytes) else str(data)
        }
        with open(self.log_file, 'a') as f:
            json.dump(log_entry, f)
            f.write('\n')  # Add a newline for readability and to separate entries

