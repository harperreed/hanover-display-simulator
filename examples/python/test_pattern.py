from pyflipdot.pyflipdot import HanoverController
from serial_proxy import SerialProxy
import config

# Create a serial port
ser = SerialProxy(config.SERIAL_PORT, config.BAUD_RATE, log_file=config.LOG_PATH)

# Create a controller
controller = HanoverController(ser)
controller.start_test_signs()
