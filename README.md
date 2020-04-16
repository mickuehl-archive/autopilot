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

export PI_KEY="~/devel/workspace/mickuehl/edgepi/bin/edgepi"
export PI_TARGET="cloudpi02"

scp -i $PI_KEY -r src cloud@$PI_TARGET:/home/cloud/

scp -i $PI_KEY -r hack cloud@$PI_TARGET:/home/cloud/

```


