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

    # test the ESC
    print("Testing the ESC")
    v.motor_speed(0.2)
    time.sleep(5)
    v.motor_speed(0)

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
    
    # test the distance sensor and acceleration sensors

    print("Testing the range finder and acceleration sensors")
    
    for i in range(20):
        print("D={:06.2f}, \tT={:04.1f}, \tacc={a[0]} {a[1]} {a[2]}, \tlinear={la[0]} {la[1]} {la[2]}, \teuler={e[0]} {e[1]} {e[2]}".format(v.distance(), v.sensor().temperature, 
            a=v.sensor().acceleration,
            la=v.sensor().linear_acceleration,
            e=v.sensor().euler,
        ))

        time.sleep(1)

    print("Done ...")
    v.shutdown()