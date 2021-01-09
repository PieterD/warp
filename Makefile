warpdev-run:
	go build -o warpdev.exe github.com/PieterD/warp/cmd/warpdev
	./warpdev.exe

docker-run: docker-build FORCE
	docker run -t -i -p 8080:8080 peted/warp-content

docker-build: clean FORCE
	cd misc/docker && sh -e ./build.sh && docker build -t peted/warp-content .

clean: FORCE
	@rm -rf ./misc/docker/target

docker-push: docker-build FORCE
	docker push peted/warp-content

prune:
	docker container prune -f
	docker image prune -f

FORCE: