
# copy hacks to the Raspi
hack :
	scp -i ${PI_KEY} -r hack/*.py cloud@${PI_TARGET}:/home/cloud/hack
	chmod +x hack/*.py

# just an example how to build go for the Raspi and upload it ...
hello: hack/golang/hello.go
	GOARM=7 GOARCH=arm GOOS=linux go build -o bin/hello hack/golang/hello.go
	scp -i ${PI_KEY} -r bin/hello cloud@${PI_TARGET}:/home/cloud/

