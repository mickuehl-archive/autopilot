#!/usr/bin/env python

import logging
import time
import argparse
import json
import sys

from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTClient
from pynput import keyboard

# init the client
mqttClient = None

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


def on_press(key):
  try:
    ch = key.char
    if ch in ['w','s','d','a','q','y','o','m','x']:
      send_message(ch)
  except AttributeError:
    return


def on_release(key):
  if key == keyboard.Key.esc:
    # Stop listener
    return False


def send_message(c):
  sys.stdout.write(c)
  # create the payload
  msg = {}
  msg['k'] = c
  msg['id'] = clientId
  msg['t'] = int(round(time.time() * 1000))
  # send the message
  mqttClient.publish(topic, json.dumps(msg), 1)


# Read in command-line parameters
parser = argparse.ArgumentParser()
parser.add_argument("-e", "--endpoint", action="store", required=True, dest="endpoint", help="Your AWS IoT custom endpoint")
parser.add_argument("-r", "--rootCA", action="store", required=True, dest="rootCAPath", help="Root CA file path")
parser.add_argument("-c", "--cert", action="store", dest="certificatePath", help="Certificate file path")
parser.add_argument("-k", "--key", action="store", dest="privateKeyPath", help="Private key file path")
parser.add_argument("-id", "--clientId", action="store", dest="clientId", default="pod-00", help="Targeted client id")
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
time.sleep(2)

print("\nReady ...\n")

# collect events until ESC
with keyboard.Listener(on_press=on_press, on_release=on_release) as listener:
  listener.join()

# cleaning up
mqttClient.disconnect()
print("\nDisconnected ...")
