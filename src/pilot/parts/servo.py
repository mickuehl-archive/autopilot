
import time
from adafruit_motor import servo


class Servo:

    def __init__(self, channel=None, range=130, trim=0):

        self.channel = channel
        self.servo_range = range / 2
        self.servo_trim = trim
        self.angle = self.servo_range

        self.servo = servo.Servo(self.channel)
        self.servo.angle = self.servo_range

        self.running = True
        print('Servo created')

    def update(self):
        while self.running:
            self.servo.angle = self.angle

    def run_threaded(self, angle):
        self.angle = (self.servo_range * angle) + self.servo_range + self.servo_trim 

    # angle = [-1.0 .. +1.0]
    def run(self, angle):
        self.run_threaded(angle)
        self.servo.angle = self.angle

    def shutdown(self):
        # set steering straight
        self.angle = self.servo_range
        self.servo.angle = self.servo_range

        time.sleep(0.3)
        self.running = False
