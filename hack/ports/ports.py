#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import time
import traceback

if __name__ == '__main__':

    data = b'aabb'

    try:
        looper = 0
        while True:
            #print("n={}".format(looper))
            #sys.stdout.buffer.write(bytes(data))
            sys.stdout.buffer.write(data)
            sys.stdout.flush()
            time.sleep(10)
            looper = looper + 1

    except KeyboardInterrupt:
        pass
    except Exception as e:
        traceback.print_exc()
    finally:
        print("Done")
