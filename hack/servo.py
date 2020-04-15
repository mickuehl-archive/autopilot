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
    for i in range(0, 90):
        v.direction(i)
        time.sleep(0.05)

    time.sleep(1)
    for i in range(0, 90):
        v.direction(i * -1)
        time.sleep(0.05)

    time.sleep(1)
    v.direction(0)
    
    print("Done ...")
    v.shutdown()