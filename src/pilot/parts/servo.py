
import time
from adafruit_motor import servo


class Servo:

    def __init__(self, channel=None, angle=130, trim=0):

        self.channel = channel
        self.servo_angle = angle
        self.servo_range = angle / 2
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
        a = (self.servo_range * angle) + self.servo_range + self.servo_trim 
        if a < 0:
            a = 0
        if a > self.servo_angle:
            a = self.servo_angle
        self.angle = a

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
