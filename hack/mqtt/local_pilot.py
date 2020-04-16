#!/usr/bin/env python

import logging
import time
import argparse
import json
import sys

from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTClient
from adafruit_motorkit import MotorKit

# init the mqtt client
mqttClient = None

# motor states
MOTOR_INC = 0.1
MOTOR_MAX = 1.0
MOTOR_MIN = -1.0

class Motors:
  def __init__(self, speed):
    self.mleft = speed
    self.mright = speed
    self.kit = MotorKit()
    self.stop()

  def stop(self):
    self.mleft = 0.0
    self.mright = 0.0
    self.set_left()
    self.set_right()

  def accelerate(self):
    self.inc_left()
    self.inc_right()

  def decelerate(self):
    self.dec_left()
    self.dec_right()

  def inc_left(self):
    self.mleft = self.mleft + MOTOR_INC
    if self.mleft > MOTOR_MAX:
      self.mleft = MOTOR_MAX
    self.set_left()

  def dec_left(self):
    self.mleft = self.mleft - MOTOR_INC
    if self.mleft < MOTOR_MIN:
      self.mleft = MOTOR_MIN
    self.set_left()

  def inc_right(self):
    self.mright = self.mright + MOTOR_INC
    if self.mright > MOTOR_MAX:
      self.mright = MOTOR_MAX
    self.set_right()

  def dec_right(self):
    self.mright = self.mright - MOTOR_INC
    if self.mright < MOTOR_MIN:
      self.mright = MOTOR_MIN
    self.set_right()

  def set_left(self):
    self.kit.motor2.throttle = self.mleft
    #print("L", self.mleft)

  def set_right(self):
    self.kit.motor1.throttle = self.mright
    #print("R", self.mright)


# lambdas
def cmd_motors_stop():
  motors.stop()

# increase the speed of the motors
def cmd_motor_left_inc():
  motors.inc_left()

def cmd_motor_right_inc():
  motors.inc_right()

# decrease the speed of the motors
def cmd_motor_left_dec():
  motors.dec_left()

def cmd_motor_right_dec():
  motors.dec_right()

def cmd_motors_acc():
  motors.accelerate()

def cmd_motors_dec():
  motors.decelerate()


# mqtt callback
def receiverCallback(client, userdata, message):
  k = json.loads(message.payload)['k']
  # look the command up
  cmd = cmd_switcher.get(k, cmd_motors_stop)
  # execute the command (it's a lambda!)
  cmd()
  

# setup
def create_logger(level):
  logger = logging.getLogger("AWSIoTPythonSDK.core")
  logger.setLevel(level)
  streamHandler = logging.StreamHandler()
  formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
  streamHandler.setFormatter(formatter)
  logger.addHandler(streamHandler)
  return logger


def create_mqtt_client(clientId, host, port, rootCAPath, privateKeyPath, certificatePath):
  mqttc = AWSIoTMQTTClient(clientId)
  # connection configuration
  mqttc.configureEndpoint(host, port)
  mqttc.configureCredentials(rootCAPath, privateKeyPath, certificatePath)
  mqttc.configureAutoReconnectBackoffTime(1, 32, 20)
  # infinite offline publish queueing
  mqttc.configureOfflinePublishQueueing(-1)
  mqttc.configureDrainingFrequency(2)  # Draining: 2 Hz
  mqttc.configureConnectDisconnectTimeout(10)  # 10 sec
  mqttc.configureMQTTOperationTimeout(30)  # 30 sec
  return mqttc


print("Booting ...")

# Read in command-line parameters
parser = argparse.ArgumentParser()
parser.add_argument("-e", "--endpoint", action="store", required=True, dest="endpoint", help="Your AWS IoT custom endpoint")
parser.add_argument("-r", "--rootCA", action="store", required=True, dest="rootCAPath", help="Root CA file path")
parser.add_argument("-c", "--cert", action="store", dest="certificatePath", help="Certificate file path")
parser.add_argument("-k", "--key", action="store", dest="privateKeyPath", help="Private key file path")
parser.add_argument("-id", "--clientId", action="store", dest="clientId", default="cloudpi00", help="Targeted client id")
parser.add_argument("-t", "--topic", action="store", dest="topic", default="cloudpi/test", help="Targeted topic")
args = parser.parse_args()

# defaults
endpoint = args.endpoint
rootCAPath = args.rootCAPath
certificatePath = args.certificatePath
privateKeyPath = args.privateKeyPath
port = 8883
clientId = args.clientId
topic = args.topic

# configure logging
logger = create_logger(logging.ERROR)

# create the MQTT client
mqttClient = create_mqtt_client(clientId, endpoint, port, rootCAPath, privateKeyPath, certificatePath)

# connect and subscribe to AWS IoT
mqttClient.connect()
mqttClient.subscribe(topic, 1, receiverCallback)
time.sleep(2)

cmd_switcher = {
  'q': cmd_motor_left_inc,
  'y': cmd_motor_left_dec,
  'o': cmd_motor_right_inc,
  'm': cmd_motor_right_dec,
  'x': cmd_motors_stop,
  'w': cmd_motors_acc,
  's': cmd_motors_dec
}
# initialize the motors
print("Starting the engines ...")
motors = Motors(0.0)

print("Go!")

try:
  while True:
    time.sleep(0.1)
except KeyboardInterrupt:
  pass

# cleaning up
mqttClient.disconnect()
print("\nDisconnected.")
