import time
from serial_proxy import SerialProxy
import config

ser = SerialProxy(config.SERIAL_PORT, config.BAUD_RATE, log_file=config.LOG_PATH)

while True:
    ser.write(config.TEST_PACKET)
    time.sleep(1)
