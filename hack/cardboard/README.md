
## Install

Copy all code to the Raspi:

```shell

scp -i ~/.ssh/cloudpi_rsa -r src pi@cloudpi<number>:/home/pi/

```

## Run the example

python src/remote_pilot.py -e a36ob0439j5iyd-ats.iot.eu-central-1.amazonaws.com -r certs/root-CA.crt -c certs/cloudpi00.cert.pem -k certs/cloudpi00.private.key



scp -i ~/.ssh/cloudpi_rsa -r pi@cloudpi02:/home/pi/data data