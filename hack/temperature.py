#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import traceback
from vehicle import Vehicle

SLEEP = 0.5

def run_tests(v):
    running = True
    print("Testing the temperature sensor")

    while running:
        print("Temperature (C):\t{}".format(v.sensor().temperature))
        time.sleep(SLEEP)


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
