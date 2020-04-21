#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from vehicle import Vehicle

if __name__ == '__main__':
    v = Vehicle()

    # test the break lights
    print("Testing LEDs")
    v.left_led_on()
    v.right_led_on()
    time.sleep(5)
    v.left_led_off()
    v.right_led_off()

    # blink
    for i in range(10):
        v.left_led_on()
        v.right_led_on()
        time.sleep(0.5)
        v.left_led_off()
        v.right_led_off()
        time.sleep(0.5)

    print("Done ...")
    v.shutdown()