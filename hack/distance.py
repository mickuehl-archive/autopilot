#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from vehicle import Vehicle

if __name__ == '__main__':
    v = Vehicle()

    # test the distance sensor and acceleration sensors

    print("Testing the range finder")
    
    for i in range(20):
        print("Measured Distance = %.1f cm" % v.distance())
        time.sleep(1)

    print("Done ...")
    v.shutdown()