import serial
import time
from serial_proxy import SerialProxy
# ser = serial.Serial('/dev/pts/3', 4800)  # Use the other port created by socat
ser = SerialProxy('/dev/pts/3', 4800, log_file="proxy.log")

# Example packet: STX + Command + Address + Resolution + Pixel Data + ETX + Checksum
test_packet = b'\x02\x11\x01\x00\xC0' + b'\xFF' * 192 + b'\x03\x00\x00'

while True:
    ser.write(test_packet)
    time.sleep(1)
