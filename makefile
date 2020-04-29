.PHONY = all clean tests deploy

PLATFORM_ARM = GOARM=7 GOARCH=arm GOOS=linux

all: tests calibrate deploy

clean:
	rm bin/calibrate bin/selftest bin/unittest

tests: unittest selftest

calibrate: cmd/calibrate/main.go
	cd cmd/calibrate && ${PLATFORM_ARM} go build -o ../../bin/calibrate main.go
	
unittest: test/unittest/main.go
	go build -o bin/unittest test/unittest/main.go

selftest: test/selftest/main.go
	${PLATFORM_ARM} go build -o bin/selftest test/selftest/main.go

deploy:
	scp -i ${PI_KEY} -r bin/calibrate cloud@${PI_TARGET}:/home/cloud/
	scp -i ${PI_KEY} -r bin/selftest cloud@${PI_TARGET}:/home/cloud/