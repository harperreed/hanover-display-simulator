import numpy as np
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial_proxy import SerialProxy

# Create a serial port (update with port name on your system)
ser = SerialProxy('/dev/pts/3', 4800, log_file="proxy.log")

# Create a controller
controller = HanoverController(ser)
controller.start_test_signs()
