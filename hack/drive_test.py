#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from vehicle import Vehicle

if __name__ == '__main__':
    v = Vehicle()

    t = 1
    s = 0.25

    # test the ESC
    print("Testing the drive for {} seconds at {} speed".format(t, s))

    v.direction(0)
    v.motor_speed(s)
    time.sleep(t)
    
    # stop
    v.motor_speed(0)

    print("Done ...")
    v.shutdown()