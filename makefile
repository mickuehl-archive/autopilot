
# copy hacks to the Raspi
hack :
	scp -i ${PI_KEY} -r hack/*.py cloud@${PI_TARGET}:/home/cloud/hack/
