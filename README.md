# autopilot
The autopilot controlling basic car functionality


## Stuff

### Setup a python virtual environment

```shell
python3 -m virtualenv -p python3 venv --system-site-packages
echo "source venv/bin/activate"
```

Either add `source venv/bin/activate` to `.bashrc` or execute it after logging in.

### Install python dependencies

```shell
pip -r src/requirements.txt
```

### Copy code to the Raspi

```shell

export PI_KEY="~/devel/workspace/mickuehl/edgepi/bin/edgepi"
export PI_TARGET="cloudpi02"

scp -i $PI_KEY -r src cloud@$PI_TARGET:/home/cloud/

scp -i $PI_KEY -r hack cloud@$PI_TARGET:/home/cloud/

```

### Makefile

