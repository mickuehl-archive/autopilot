.PHONY = all test deploy

PLATFORM = GOARM=7 GOARCH=arm GOOS=linux

all: test calibrate deploy

calibrate: cmd/calibrate/main.go
	${PLATFORM} go build -o bin/calibrate cmd/calibrate/main.go

test: test/selftest.go pkg/pilot/pilot.go pkg/pilot/raspi.go
	${PLATFORM} go build -o bin/selftest test/selftest.go

deploy:
	scp -i ${PI_KEY} -r bin/* cloud@${PI_TARGET}:/home/cloud/