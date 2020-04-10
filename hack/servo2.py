#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import busio, time
from board import SCL, SDA
from adafruit_pca9685 import PCA9685
from adafruit_motor import servo

CHANNEL = 0
FREQ = 50

def set_servo(servo, a):
    print("a=", a)
    servo.angle = a
    time.sleep(2)

i2c = busio.I2C(SCL, SDA)

# Create a simple PCA9685 class instance.
pca = PCA9685(i2c)
pca.frequency = FREQ

servo1 = servo.Servo(pca.channels[CHANNEL])

time.sleep(5)
set_servo(servo1, 0)
set_servo(servo1, 45)
set_servo(servo1, 90)
set_servo(servo1, 135)
set_servo(servo1, 180)
set_servo(servo1, 0)

pca.deinit()
