.PHONY = all clean tests deploy hacks frameserver picam

PLATFORM_ARM = GOARM=7 GOARCH=arm GOOS=linux

all: tests deploy

clean:
	rm bin/unittest

tests: unittest frameserver
	
unittest: test/unittest/main.go
	go build -o bin/unittest test/unittest/main.go

frameserver: frameserver/main.go
	${PLATFORM_ARM} go build -o ../bin/frameserver frameserver/main.go
	scp -i ${PI_KEY} -r ../bin/frameserver cloud@${PI_TARGET}:/home/cloud/

picam: picam/main.go
	${PLATFORM_ARM} go build -o ../bin/picam picam/main.go
	scp -i ${PI_KEY} -r ../bin/picam cloud@${PI_TARGET}:/home/cloud/

hacks:
	${PLATFORM_ARM} go build -o ../bin/frameserver frameserver/main.go
	scp -i ${PI_KEY} -r ../bin/frameserver cloud@${PI_TARGET}:/home/cloud/
	#${PLATFORM_ARM} go build -o bin/picam test/picam/main.go
	#scp -i ${PI_KEY} -r bin/picam cloud@${PI_TARGET}:/home/cloud/
	#scp -i ${PI_KEY} -r pkg/parts/camera/camera.py cloud@${PI_TARGET}:/home/cloud/

deploy:
	scp -i ${PI_KEY} -r bin/selftest cloud@${PI_TARGET}:/home/cloud/
