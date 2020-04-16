#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from cardboard import Vehicle, Motor
from cardboard.parts.datastore import TubWriter
from cardboard.parts.camera import PiCamera
from cardboard.parts.driver import Driver
from cardboard.parts.controller import LocalWebController

IMAGE_W = 640
IMAGE_H = 480
IMAGE_DEPTH = 1


if __name__ == '__main__':
    v = Vehicle()

    cam = PiCamera(image_w=IMAGE_W, image_h=IMAGE_H, image_d=IMAGE_DEPTH, vflip=True, hflip=True)
    v.add(cam, outputs=['cam/image'], threaded=True)

    v.add(LocalWebController(), inputs=['cam/image'],
          outputs=['user/angle', 'user/throttle', 'user/mode', 'recording'], threaded=True)

    # the drive-train
    v.add(Driver(), inputs=['user/throttle','user/angle'], outputs=['motor/left', 'motor/right'], threaded=True)

    #tub = TubWriter(path='~/data', inputs=['image','s'], types=['image_array','float'])
    #v.add(tub, inputs=['image','s'], outputs=['num_records'])

    # start the drive loop at 10 Hz for 3 sec.
    #v.start(rate_hz=10, max_loop_count=30, verbose=True)
    v.start(rate_hz=10)