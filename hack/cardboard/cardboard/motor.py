#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from adafruit_motorkit import MotorKit

# motor states
MOTOR_INC = 0.1
MOTOR_MAX = 1.0
MOTOR_MIN = -1.0
MOTOR_THRESHOLD = 0.9
MOTOR_DELAY = 0.5

class Motor:
  def __init__(self, aleft=1.0, aright=1.0):
    self.mleft = 0.0 # speed
    self.mright = 0.0
    self.aleft = aleft # motor trim
    self.aright = aright
    self.kit = MotorKit()
    self.stop()

  def stop(self):
    self.set_left(0.0)
    self.set_right(0.0)

  def accelerate(self):
    self.inc_left()
    self.inc_right()

  def decelerate(self):
    self.dec_left()
    self.dec_right()

  def forward(self, speed):
    if speed > MOTOR_MAX:      
      speed = MOTOR_MAX

    # overcome the initial resistance when starting from a stop
    if self.mright == 0 and self.mleft == 0:
      self.kit.motor1.throttle = MOTOR_THRESHOLD
      self.kit.motor2.throttle = MOTOR_THRESHOLD
      time.sleep(MOTOR_DELAY)

    self.set_left(speed)
    self.set_right(speed)
    
  def inc_left(self, inc = MOTOR_INC):
    self.set_left(self.mleft + inc)

  def dec_left(self, dec = MOTOR_INC):
    self.set_left(self.mleft - dec)

  def inc_right(self, inc = MOTOR_INC):
    self.set_right(self.mright + inc)

  def dec_right(self, dec = MOTOR_INC):
    self.set_right(self.mright - dec)

  def set_left(self, s):
    if self.mleft == s:
      return
    if s < MOTOR_MIN:
      self.mleft = 0.0
    elif s > MOTOR_MAX:
      self.mleft = MOTOR_MAX
    else:
      self.mleft = s

    self.kit.motor2.throttle = self.mleft * self.aleft

  def set_right(self, s):
    if self.mright == s:
      return
    if s < MOTOR_MIN:
      self.mright = 0.0
    elif s > MOTOR_MAX:
      self.mright = MOTOR_MAX
    else:
      self.mright = s

    self.kit.motor1.throttle = self.mright * self.aright