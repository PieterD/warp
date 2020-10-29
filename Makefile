
build-docker: FORCE
	rm -rf misc/docker/target
	mkdir -p misc/docker/target
	cp -r misc/static/* misc/docker/target

FORCE: