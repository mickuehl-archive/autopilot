#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import busio

from board import SCL, SDA
from adafruit_pca9685 import PCA9685
from adafruit_motor import servo
import adafruit_bno055
import RPi.GPIO as GPIO

PCA_FREQUENCY = 50

LEFT_LED_CHANNEL = 11
RIGHT_LED_CHANNEL = 8

SERVO_CHANNEL = 3
ABSOLUTE_SERVO_RANGE = 120
SERVO_TRIM = -4

LED_OFF = 0
LED_ON = 0xFFFF

ESC_CHANNEL = 0
FORWARD = 1.0
REVERSE = -1.0
ESC_MAX_SPEED = 0.25

# set GPIO Pins for the forward distance sensor
GPIO_TRIGGER = 23
GPIO_ECHO = 24


class Vehicle:

    def __init__(self):

        # GPIO Mode (BOARD / BCM)
        GPIO.setmode(GPIO.BCM)
        # set GPIO direction (IN / OUT)
        GPIO.setup(GPIO_TRIGGER, GPIO.OUT)
        GPIO.setup(GPIO_ECHO, GPIO.IN)

        # Create the I2C bus interface.
        i2c_bus = busio.I2C(SCL, SDA)

        # Create a simple PCA9685 class instance.
        self.pca = PCA9685(i2c_bus)

        # Set the PWM frequency to e.g. 50hz.
        self.pca.frequency = PCA_FREQUENCY

        # the different actors/sensors

        # break lights
        self.led1 = self.pca.channels[LEFT_LED_CHANNEL]
        self.led2 = self.pca.channels[RIGHT_LED_CHANNEL]
        self.led1.duty_cycle = LED_OFF
        self.led2.duty_cycle = LED_OFF

        # ESC
        self.esc = servo.ContinuousServo(self.pca.channels[ESC_CHANNEL])
        self.max_speed = ESC_MAX_SPEED
        self.motor_speed(0)

        # servo
        self.servo1 = servo.Servo(self.pca.channels[SERVO_CHANNEL])
        self.servo_abs_range = ABSOLUTE_SERVO_RANGE
        self.servo_range = self.servo_abs_range / 2
        self.servo_trim = SERVO_TRIM
        self.direction(0)

        # Bosch BNO055 sensor
        self.bosch = adafruit_bno055.BNO055(i2c_bus)

    def shutdown(self):
        self.pca.deinit()
        
    def left_led_on(self):
        self.led1.duty_cycle = LED_ON

    def left_led_off(self):
        self.led1.duty_cycle = LED_OFF

    def right_led_on(self):
        self.led2.duty_cycle = LED_ON

    def right_led_off(self):
        self.led2.duty_cycle = LED_OFF

    def motor_speed(self, speed, direction=FORWARD):
        if speed > self.max_speed:
            speed = self.max_speed
        self.esc.throttle = speed * direction

    def direction(self, angle):
        if angle > self.servo_range:
            angle = self.servo_range
        if angle < self.servo_range * -1:
            angle = self.servo_range * -1

        self.servo1.angle = self.servo_range +self.servo_trim + angle

    def sensor(self):
        return self.bosch

    def distance(self):
        try:
            # set Trigger to HIGH
            GPIO.output(GPIO_TRIGGER, True)

            # set Trigger after 0.01ms to LOW
            time.sleep(0.00001)
            GPIO.output(GPIO_TRIGGER, False)

            StartTime = time.time()
            StopTime = time.time()

            # save StartTime
            while GPIO.input(GPIO_ECHO) == 0:
                StartTime = time.time()

            # save time of arrival
            while GPIO.input(GPIO_ECHO) == 1:
                StopTime = time.time()

            # time difference between start and arrival
            TimeElapsed = StopTime - StartTime
            # multiply with the sonic speed (34300 cm/s)
            # and divide by 2, because there and back
            distance = (TimeElapsed * 34300) / 2

            return distance

        except:
            return -1.0
