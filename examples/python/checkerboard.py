import numpy as np
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial_proxy import SerialProxy
import config

# Create a serial port
print("DEBUG: Creating serial port with config: PORT={}, BAUD={}".format(config.SERIAL_PORT, config.BAUD_RATE))
ser = SerialProxy(config.SERIAL_PORT, config.BAUD_RATE, log_file=config.LOG_PATH)
print("DEBUG: Serial port created successfully")

# Create a controller
print("DEBUG: Creating HanoverController")
controller = HanoverController(ser)
print("DEBUG: HanoverController created successfully")

# Add a sign
print("DEBUG: Creating HanoverSign with address={}, width={}, height={}".format(config.SIGN_ADDRESS, config.SIGN_WIDTH, config.SIGN_HEIGHT))
sign = HanoverSign(address=config.SIGN_ADDRESS, width=config.SIGN_WIDTH, height=config.SIGN_HEIGHT)
print("DEBUG: HanoverSign created successfully")
print("DEBUG: Adding sign to controller with name 'dev'")
controller.add_sign('dev', sign)
print("DEBUG: Sign added to controller successfully")

# Create a 'checkerboard' image
print("DEBUG: Creating checkerboard image")
image = sign.create_image()
print("DEBUG: Image created with shape:", image.shape)
print("DEBUG: Setting checkerboard pattern with interval:", config.CHECKERBOARD_INTERVAL)
image[::config.CHECKERBOARD_INTERVAL, ::config.CHECKERBOARD_INTERVAL] = True
image[1::config.CHECKERBOARD_INTERVAL, 1::config.CHECKERBOARD_INTERVAL] = True
print("DEBUG: Checkerboard pattern set successfully")
print("DEBUG: Drawing image on controller")
controller.draw_image(image)
print("DEBUG: Image drawn successfully")
