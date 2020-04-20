hack :
	scp -i ${PI_KEY} -r hack cloud@${PI_TARGET}:/home/cloud/

hello: hack/hello.go
	GOARM=7 GOARCH=arm GOOS=linux go build -o bin/hello hack/hello.go
	scp -i ${PI_KEY} -r bin/hello cloud@${PI_TARGET}:/home/cloud/

