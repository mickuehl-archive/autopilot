#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
from vehicle import Vehicle

if __name__ == '__main__':
    v = Vehicle()

    print("Testing the acceleration sensors")
    
    for i in range(20):
        print("T={:04.1f}, \tacc={a[0]} {a[1]} {a[2]}, \tlinear={la[0]} {la[1]} {la[2]}, \teuler={e[0]} {e[1]} {e[2]}, \tgyro={g[0]} {g[1]} {g[2]}, \tmag={m[0]} {m[1]} {m[2]}, \tqat={q[0]} {q[1]} {q[2]}".format(v.sensor().temperature, 
            a=v.sensor().acceleration,
            m=v.sensor().magnetic,
            g=v.sensor().gyro,
            e=v.sensor().euler,
            la=v.sensor().linear_acceleration,
            q=v.sensor().quaternion
        ))

        time.sleep(1)

    print()
    print("Done ...")
    v.shutdown()