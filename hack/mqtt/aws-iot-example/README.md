# AWS IoT Device SDK

## Installation

### Certificates

Install the certificates on the device

```shell
export DEVICE_NAME=cloudpi01

CERTS="certs/cloudpi00.cert.pem certs/cloudpi00.private.key certs/root-CA.crt"
scp $CERTS pi@$DEVICE_NAME:/home/pi/certs/

```
### Source code

Upload the source code

```shell

scp -r src pi@$DEVICE_NAME:/home/pi

```

## Run the example

python src/main.py -e a36ob0439j5iyd-ats.iot.eu-central-1.amazonaws.com -r certs/root-CA.crt -c certs/cloudpi01.cert.pem -k certs/cloudpi01.private.key

python src/main.py -e a36ob0439j5iyd-ats.iot.eu-central-1.amazonaws.com -r ../../certs/root-CA.crt -c ../../certs/cloudpi00.cert.pem -k ../../certs/cloudpi00.private.key

python src/main.py -e a36ob0439j5iyd-ats.iot.eu-central-1.amazonaws.com -r root-CA.crt -c cloudpi.cert.pem -k cloudpi.private.key

python src/main.py -e a36ob0439j5iyd-ats.iot.eu-central-1.amazonaws.com -r certs/ca/root-CA.crt -c certs/cert/cloudpi.cert.pem -k certs/key/cloudpi.private.key


## Deploy to OpenShift

oc new-project aws-iot-example

oc create configmap root-ca --from-file=certs/root-CA.crt
oc create configmap cloudpi-cert --from-file=certs/cloudpi00.cert.pem
oc create configmap cloudpi-key --from-file=certs/cloudpi00.private.key

oc create -f examples/aws-iot-example/deploy.yaml



### References

* https://github.com/aws/aws-iot-device-sdk-python

#### Time-out issues

* https://github.com/aws/aws-iot-device-sdk-python/issues/154




oc create configmap <configmap-name> --from-file=</path/to/file/ca.pem>
