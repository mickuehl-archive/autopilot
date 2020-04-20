#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import traceback
from vehicle import Vehicle


def run_tests(v):

    # test the ESC
    print("Testing the ESC")
    v.motor_speed(0.2)
    time.sleep(5)
    v.motor_speed(0)

    # test the servo
    print("Testing the servo")
    v.direction(0)
    time.sleep(5)

    # wiggle left/right a bit
    v.direction(30)
    time.sleep(0.4)

    v.direction(-30)
    time.sleep(0.4)

    v.direction(0)


if __name__ == '__main__':
    try:
        v = Vehicle()
        run_tests(v)
    except KeyboardInterrupt:
        pass
    except Exception as e:
        traceback.print_exc()
    finally:
        v.shutdown()

    print()
    print("Done ...")
