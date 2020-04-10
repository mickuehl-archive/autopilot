#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import busio
import time
from board import SCL, SDA
from adafruit_pca9685 import PCA9685
from adafruit_motor import servo

CHANNEL = 8
FREQ = 50
MAX_SPEED = 0.25
FORWARD = 1.0
REVERSE = -1.0

i2c = busio.I2C(SCL, SDA)

# Create a simple PCA9685 class instance.
pca = PCA9685(i2c)
pca.frequency = FREQ

esc = servo.ContinuousServo(pca.channels[CHANNEL])
esc.throttle = 0

throttle = 0.0
while throttle < MAX_SPEED:
    esc.throttle = throttle
    throttle += 0.01
    time.sleep(0.2)

esc.throttle = 0
time.sleep(3)

throttle = 0.0
while throttle < MAX_SPEED:
    esc.throttle = throttle * REVERSE
    throttle += 0.01
    time.sleep(0.2)

esc.throttle = 0
pca.deinit()
