.PHONY = all clean tests deploy hacks

PLATFORM_ARM = GOARM=7 GOARCH=arm GOOS=linux

all: tests calibrate deploy

clean:
	rm bin/calibrate bin/selftest bin/unittest

tests: unittest selftest scenario1 frameserver

calibrate: cmd/calibrate/main.go
	cd cmd/calibrate && ${PLATFORM_ARM} go build -o ../../bin/calibrate main.go
	
unittest: test/unittest/main.go
	go build -o bin/unittest test/unittest/main.go

hacks:
	#${PLATFORM_ARM} go build -o bin/frameserver test/frameserver/main.go
	#scp -i ${PI_KEY} -r bin/frameserver cloud@${PI_TARGET}:/home/cloud/
	#${PLATFORM_ARM} go build -o bin/picam test/picam/main.go
	#scp -i ${PI_KEY} -r bin/picam cloud@${PI_TARGET}:/home/cloud/
	scp -i ${PI_KEY} -r pkg/parts/camera/camera.py cloud@${PI_TARGET}:/home/cloud/

selftest: test/selftest/main.go
	${PLATFORM_ARM} go build -o bin/selftest test/selftest/main.go

scenario1: test/scenario1/main.go
	${PLATFORM_ARM} go build -o bin/scenario1 test/scenario1/main.go

deploy:
	scp -i ${PI_KEY} -r bin/calibrate cloud@${PI_TARGET}:/home/cloud/
	scp -i ${PI_KEY} -r bin/selftest cloud@${PI_TARGET}:/home/cloud/
	scp -i ${PI_KEY} -r bin/scenario1 cloud@${PI_TARGET}:/home/cloud/
	