
# copy the hacks to the Raspi
hack :
	scp -i ${PI_KEY} -r *.py cloud@${PI_TARGET}:/home/cloud/hack/
	scp -i ${PI_KEY} -r *.txt cloud@${PI_TARGET}:/home/cloud/hack/
	ssh -i ${PI_KEY} cloud@${PI_TARGET} "chmod +x hack/*.py"
