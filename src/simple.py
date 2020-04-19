#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import traceback
from pilot import Pilot
from pilot.parts.camera import PiCamera
from pilot.parts.controller import LocalWebController

IMAGE_W = 640
IMAGE_H = 480
IMAGE_DEPTH = 1

if __name__ == '__main__':
    try:
        p = Pilot()

        cam = PiCamera(image_w=IMAGE_W, image_h=IMAGE_H, image_d=IMAGE_DEPTH, vflip=True, hflip=True)
        p.add(cam, outputs=['cam/image'], threaded=True)

        p.add(LocalWebController(), inputs=['cam/image'], outputs=['user/angle', 'user/throttle', 'user/mode', 'recording'], threaded=True)

        p.start()
    except KeyboardInterrupt:
        pass
    except Exception as e:
        traceback.print_exc()
    finally:
        p.stop()

    print()
    print("Done ...")
