#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import traceback
from pilot import Pilot
from pilot.parts.camera import PiCamera
from pilot.parts.web import WebController
from pilot.parts.channel import Channel
from pilot.parts.servo import Servo

IMAGE_W = 640
IMAGE_H = 480
IMAGE_DEPTH = 1

PCA9685_I2C_ADDR = 0x40

STEERING_CHANNEL = 3
STEERING_LEFT_PWM = 750
STEERING_RIGHT_PWM = 2250

if __name__ == '__main__':
    try:
        # the vehicle
        p = Pilot()

        # the Pi camera used to record the drive
        cam = PiCamera(image_w=IMAGE_W, image_h=IMAGE_H,
                       image_d=IMAGE_DEPTH, vflip=True, hflip=True)
        p.add(cam, outputs=['cam/image'], threaded=True)

        # the controller used to provide a simple interface for training
        web = WebController()
        p.add(web, inputs=['cam/image'], outputs=['user/angle',
                                                  'user/throttle', 'user/mode', 'recording'], threaded=True)

        # the car parts
        servo_channel = Channel(channel=STEERING_CHANNEL)
        servo = Servo(channel=servo_channel.channel(), angle=130, trim=-4)
        p.add(servo, inputs=['user/angle'])

        # start the vehicle's main loop
        p.start()
    except KeyboardInterrupt:
        pass
    except Exception as e:
        traceback.print_exc()
    finally:
        p.stop()

    print()
    print("Done ...")
