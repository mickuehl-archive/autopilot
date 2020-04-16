#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from vehicle import Vehicle

if __name__ == '__main__':
    v = Vehicle()

    # test the servo
    print("Testing the servo")
    v.direction(0)
    time.sleep(1)

    # wiggle left/right a bit
    v.direction(20)
    time.sleep(0.2)

    v.direction(-20)
    time.sleep(0.2)

    v.direction(0)

    v.direction(40)
    time.sleep(1)

    v.direction(0)
    time.sleep(1)

    v.direction(-40)
    time.sleep(1)

    v.direction(0)
    
    print("Done ...")
    v.shutdown()