import os

# Serial port configuration
SERIAL_PORT = '/dev/pts/1'  # Update this to the correct port for your system
BAUD_RATE = 4800

# Logging configuration
LOG_FILE = 'proxy.log'

# Display configuration
SIGN_ADDRESS = 1
SIGN_WIDTH = 96
SIGN_HEIGHT = 16

# Paths
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
LOG_PATH = os.path.join(BASE_DIR, LOG_FILE)

# Test packet configuration
TEST_PACKET = b'\x02\x11\x01\x00\xC0' + b'\xFF' * 192 + b'\x03\x00\x00'

# Checkerboard configuration
CHECKERBOARD_INTERVAL = 3  # Interval for checkerboard pattern
