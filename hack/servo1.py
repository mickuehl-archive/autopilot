#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import busio, time
from board import SCL, SDA
from adafruit_pca9685 import PCA9685
from adafruit_motor import servo

CHANNEL = 15
FREQ = 50

i2c = busio.I2C(SCL, SDA)

# Create a simple PCA9685 class instance.
pca = PCA9685(i2c)
pca.frequency = FREQ

servo1 = servo.Servo(pca.channels[CHANNEL])

for i in range(0, 180,5):
    servo1.angle = i
    time.sleep(0.3)
for i in range(180):
    servo1.angle = 180 - i

pca.deinit()
