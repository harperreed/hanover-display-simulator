import numpy as np
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial_proxy import SerialProxy

# Create a serial port (update with port name on your system)
ser = SerialProxy('/dev/pts/3', 4800, log_file="proxy.log")

# Create a controller
controller = HanoverController(ser)

# Add a sign
# Note: The sign's address i
sign = HanoverSign(address=1, width=96, height=16)
controller.add_sign('dev', sign)

# Create a 'checkerboard' image
image = sign.create_image()
image[::2, ::2] = True
image[1::2, 1::2] = True
controller.draw_image(image)
