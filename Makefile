
run-docker: build-docker FORCE
	docker run -t -i -p 8080:8080 peted/warp-content

build-docker: clean FORCE
	cd misc/docker && sh -e ./build.sh && docker build -t peted/warp-content .

clean: FORCE
	@rm -rf ./misc/docker/target

push-docker: build-docker FORCE
	docker push peted/warp-content

FORCE: