#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import sys
from cardboard.motors import Motors


def drive(speed, t, al, ar):
    motors = Motors(al, ar)

    motors.forward(speed)
    time.sleep(t)

    # cleaning up
    motors.stop()


if __name__ == '__main__':
    # initialize the motors
    print("Starting the engines ...")

    # defaults
    speed = 0.5
    t = 1.0
    al = 1.0
    ar = 1.0

    params = len(sys.argv)
    if params == 2:
        speed = float(sys.argv[1])
    elif params == 3:
        speed = float(sys.argv[1])
        t = float(sys.argv[2])
    elif params > 3:
        speed = float(sys.argv[1])
        t = float(sys.argv[2])
        al = float(sys.argv[3])
        ar = float(sys.argv[4])

    # Go !
    print("Go!")
    drive(speed, t, al, ar)

    # done
    print("\nParked.")
