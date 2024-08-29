import numpy as np
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial_proxy import SerialProxy
import config

# Create a serial port
ser = SerialProxy(config.SERIAL_PORT, config.BAUD_RATE, log_file=config.LOG_PATH)

# Create a controller
controller = HanoverController(ser)

# Add a sign
sign = HanoverSign(address=config.SIGN_ADDRESS, width=config.SIGN_WIDTH, height=config.SIGN_HEIGHT)
controller.add_sign('dev', sign)

# Create a 'checkerboard' image
image = sign.create_image()
image[::config.CHECKERBOARD_INTERVAL, ::config.CHECKERBOARD_INTERVAL] = True
image[1::config.CHECKERBOARD_INTERVAL, 1::config.CHECKERBOARD_INTERVAL] = True
controller.draw_image(image)
