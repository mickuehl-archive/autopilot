
servo: servo.go
	GOARM=7 GOARCH=arm GOOS=linux go build -o ../../bin/servo servo.go
	scp -i ${PI_KEY} -r ../../bin/servo cloud@${PI_TARGET}:/home/cloud/
