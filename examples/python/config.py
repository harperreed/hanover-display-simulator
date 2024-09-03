import os
import yaml

with open('../../config.yaml', 'r') as config_file:
    config = yaml.safe_load(config_file)

SERIAL_PORT = config['serial_port_in']  # Load from config.yaml
BAUD_RATE = config['baud_rate']

# Logging configuration
LOG_FILE = 'proxy.log'

# Display configuration
SIGN_ADDRESS = config['address']
SIGN_WIDTH = config['columns']
SIGN_HEIGHT = config['rows']

# Paths
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
LOG_PATH = os.path.join(BASE_DIR, LOG_FILE)

# Test packet configuration
TEST_PACKET = b'\x02\x11\x01\x00\xC0' + b'\xFF' * 192 + b'\x03\x00\x00'

# Checkerboard configuration
CHECKERBOARD_INTERVAL = 3  # Interval for checkerboard pattern
