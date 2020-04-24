.PHONY = all test1 deploy

PLATFORM = GOARM=7 GOARCH=arm GOOS=linux

all: test1 calibrate deploy

calibrate: cmd/calibrate/main.go
	${PLATFORM} go build -o bin/calibrate cmd/calibrate/main.go

test1: test/test1.go
	${PLATFORM} go build -o bin/test1 test/test1.go

deploy:
	scp -i ${PI_KEY} -r bin/* cloud@${PI_TARGET}:/home/cloud/