
gobot: gobot.go
	GOARM=7 GOARCH=arm GOOS=linux go build -o ../../bin/gobot gobot.go
	scp -i ${PI_KEY} -r ../../bin/gobot cloud@${PI_TARGET}:/home/cloud/
