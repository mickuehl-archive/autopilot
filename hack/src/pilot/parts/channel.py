
import time

import busio

from board import SCL, SDA
from adafruit_pca9685 import PCA9685


class Channel:

    def __init__(self, channel=0, frequency=50, init_delay=0.1):

        self.chan = channel

        # Create the I2C bus interface.
        i2c_bus = busio.I2C(SCL, SDA)
        self.pca = PCA9685(i2c_bus)
        self.pca.frequency = frequency

    def channel(self):
        return self.pca.channels[self.chan]
