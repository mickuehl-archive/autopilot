#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import traceback
from threading import Thread

from .memory import Memory


class Pilot:
    def __init__(self):
        self.on = False
        self.parts = []
        self.threads = []
        self.mem = Memory()

    def add(self, part, inputs=[], outputs=[], threaded=False, run_condition=None):

        assert type(inputs) is list, "inputs is not a list: %r" % inputs
        assert type(outputs) is list, "outputs is not a list: %r" % outputs
        assert type(
            threaded) is bool, "threaded is not a boolean: %r" % threaded

        p = part
        print('Adding part {}.'.format(p.__class__.__name__))
        entry = {}
        entry['part'] = p
        entry['inputs'] = inputs
        entry['outputs'] = outputs
        entry['run_condition'] = run_condition

        if threaded:
            t = Thread(target=part.update, args=())
            t.daemon = True
            entry['thread'] = t

        self.parts.append(entry)

    def remove(self, part):
        self.parts.remove(part)

    def update_parts(self):
        for entry in self.parts:

            run = True
            # check run condition, if it exists
            if entry.get('run_condition'):
                run_condition = entry.get('run_condition')
                run = self.mem.get([run_condition])[0]

            if run:
                # get part
                p = entry['part']

                # get inputs from memory
                inputs = self.mem.get(entry['inputs'])
                # run the part
                if entry.get('thread'):
                    outputs = p.run_threaded(*inputs)
                else:
                    outputs = p.run(*inputs)

                # save the output to memory
                if outputs is not None:
                    self.mem.put(entry['outputs'], outputs)

    def start(self, rate_hz=10, max_loop_count=None, verbose=False):
        try:

            self.on = True

            for entry in self.parts:
                if entry.get('thread'):
                    entry.get('thread').start()

            # wait until the parts warm up.
            print('Starting pilot at {} Hz'.format(rate_hz))

            loop_count = 0
            while self.on:
                start_time = time.time()
                loop_count += 1

                self.update_parts()

                # stop drive loop if loop_count exceeds max_loopcount
                if max_loop_count and loop_count > max_loop_count:
                    self.on = False

                sleep_time = 1.0 / rate_hz - (time.time() - start_time)
                if sleep_time > 0.0:
                    time.sleep(sleep_time)
                else:
                    # print a message when could not maintain loop rate.
                    if verbose:
                        print('WARN::Pilot: jitter violation in vehicle loop '
                              'with {0:4.0f}ms'.format(abs(1000 * sleep_time)))

        except KeyboardInterrupt:
            pass
        except Exception as e:
            traceback.print_exc()
        finally:
            self.stop()

    def stop(self):
        print('Shutting down pilot and its parts...')
        self.on = False

        # for entry in self.parts:
        #    try:
        #        entry['part'].shutdown()
        #    except AttributeError:
        #        # usually from missing shutdown method, which should be optional
        #        pass
        #    except Exception as e:
        #        print(e)
