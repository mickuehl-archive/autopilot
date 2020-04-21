#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from vehicle import Vehicle

if __name__ == '__main__':
    v = Vehicle()

    # test the ESC
    print("Testing the ESC")

    print("Speed = 0.2")
    v.motor_speed(0.2)
    time.sleep(1)
    
    print("Speed = 0.1")
    v.motor_speed(0.1)
    time.sleep(2)

    # accelerate and break
    print("Speed = 0.2")
    v.motor_speed(0.2)
    time.sleep(1)
    print("Break = -0.1")
    v.motor_speed(0.1, -1.0)
    time.sleep(1)

    # reverse
    print("Speed = -0.2")
    v.motor_speed(0.2, -1.0)
    time.sleep(2)

    # stop
    v.motor_speed(0)

    print("Done ...")
    v.shutdown()