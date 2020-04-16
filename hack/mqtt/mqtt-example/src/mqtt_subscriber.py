#!/usr/bin/env python

import paho.mqtt.client as mqtt
import argparse


def on_connect(client, userdata, flags, rc):
    print("Connected with result code "+str(rc))


def on_message(client, userdata, msg):
    print("topic -> %s %s" % (msg.topic, msg.payload))


def on_publish(mqttc, obj, mid):
    print("mid: " + str(mid))


def on_subscribe(mqttc, obj, mid, granted_qos):
    print("Subscribed: " + str(mid) + " " + str(granted_qos))


# args
parser = argparse.ArgumentParser()
parser.add_argument("-e", "--endpoint", action="store", required=True,
                    dest="endpoint", help="Your AWS IoT custom endpoint")
parser.add_argument("-r", "--rootCA", action="store",
                    required=True, dest="rootCAPath", help="Root CA file path")
parser.add_argument("-c", "--cert", action="store",
                    dest="certificatePath", help="Certificate file path")
parser.add_argument("-k", "--key", action="store",
                    dest="privateKeyPath", help="Private key file path")
parser.add_argument("-id", "--clientId", action="store", dest="clientId",
                    default="pod-mqtt-receiver", help="Targeted client id")
parser.add_argument("-t", "--topic", action="store", dest="topic",
                    default="cloudpi/test", help="Targeted topic")
parser.add_argument("-p", "--port", action="store",
                    dest="port", type=int, help="Port number override")
args = parser.parse_args()

endpoint = args.endpoint
rootCAPath = args.rootCAPath
certificatePath = args.certificatePath
privateKeyPath = args.privateKeyPath
clientId = args.clientId
topic = args.topic
port = 8883
if args.port:
    port = args.port

client = mqtt.Client(client_id=clientId)

client.tls_set(ca_certs=rootCAPath, certfile=certificatePath,
               keyfile=privateKeyPath)
client.on_connect = on_connect
client.on_message = on_message
client.on_publish = on_publish
client.on_subscribe = on_subscribe

# connect and subscribe
client.connect(endpoint, port, 60)
client.subscribe(topic)

# wait for messages
client.loop_forever()
