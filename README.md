# autopilot
The autopilot controlling basic car functionality


## Stuff

### Setup a python virtual env

```shell
python3 -m virtualenv -p python3 venv --system-site-packages
echo "source venv/bin/activate"
```

### Copy code to the Raspberry Pi

```shell

export SSH_KEY="~/devel/workspace/mickuehl/edgepi/bin/edgepi"
export TARGET_PI="cloudpi02"

scp -i $SSH_KEY -r src cloud@$TARGET_PI:/home/cloud/

scp -i $SSH_KEY -r hack cloud@$TARGET_PI:/home/cloud/

```



https://circuitpython.readthedocs.io/projects/pca9685/en/latest/index.html
https://circuitpython.readthedocs.io/projects/servokit/en/latest/

https://github.com/adafruit/Adafruit_CircuitPython_Motor
https://github.com/adafruit/Adafruit_CircuitPython_ServoKit
