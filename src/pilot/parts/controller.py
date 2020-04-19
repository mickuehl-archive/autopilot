
import os
import json
import time
import asyncio

import requests

from socket import gethostname

from tornado.ioloop import IOLoop
from tornado.web import Application, RedirectHandler, StaticFileHandler, RequestHandler
from tornado.httpserver import HTTPServer
import tornado.gen

from .. import utils


class LocalWebController(tornado.web.Application):

    def __init__(self):
        print('Starting Server...', end='')

        this_dir = os.path.dirname(os.path.realpath(__file__))
        self.static_file_path = os.path.join(this_dir, 'templates', 'static')
        self.angle = 0.0
        self.throttle = 0.0
        self.mode = 'user'
        self.recording = False
        self.port = 8887

        handlers = [
            (r"/", RedirectHandler, dict(url="/drive")),
            (r"/drive", DriveAPI),
            (r"/video", VideoAPI),
            (r"/static/(.*)", StaticFileHandler,
             {"path": self.static_file_path}),
        ]

        settings = {'debug': True}
        super().__init__(handlers, **settings)

        print("... you can now go to {}.local:8887 to drive your car.".format(
            gethostname()))

    def update(self):
        asyncio.set_event_loop(asyncio.new_event_loop())
        self.listen(self.port)
        IOLoop.instance().start()

    def run_threaded(self, img_arr=None):
        self.img_arr = img_arr
        return self.angle, self.throttle, self.mode, self.recording

    def run(self, img_arr=None):
        self.img_arr = img_arr
        return self.angle, self.throttle, self.mode, self.recording

    def shutdown(self):
        pass


class DriveAPI(RequestHandler):

    def get(self):
        data = {}
        self.render("templates/vehicle.html", **data)

    def post(self):
        data = tornado.escape.json_decode(self.request.body)
        self.application.angle = data['angle']
        self.application.throttle = data['throttle']
        self.application.mode = data['drive_mode']
        self.application.recording = data['recording']


class VideoAPI(RequestHandler):

    async def get(self):

        self.set_header(
            "Content-type", "multipart/x-mixed-replace;boundary=--boundarydonotcross")

        served_image_timestamp = time.time()
        my_boundary = "--boundarydonotcross\n"
        while True:

            interval = .01
            if served_image_timestamp + interval < time.time() and hasattr(self.application, 'img_arr'):

                img = utils.arr_to_binary(self.application.img_arr)
                self.write(my_boundary)
                self.write("Content-type: image/jpeg\r\n")
                self.write("Content-length: %s\r\n\r\n" % len(img))
                self.write(img)
                served_image_timestamp = time.time()
                try:
                    await self.flush()
                except tornado.iostream.StreamClosedError:
                    pass
            else:
                await tornado.gen.sleep(interval)
