warpdev-run: FORCE
	go build -o warpdev.exe github.com/PieterD/warp/cmd/warpdev
	./warpdev.exe


app-prog-build: FORCE
	cd misc/docker && sh -e ./build.sh prog "-manifest"

app-prog-docker: app-prog-build FORCE
	cd misc/docker && docker build -t peted/warp-content .

app-prog-run: app-prog-docker FORCE
	docker run -t -i -p 8080:8080 peted/warp-content

app-gltest-build: FORCE
	cd misc/docker && sh -e ./build.sh gltest

app-gltest-docker: app-gltest-build FORCE
	cd misc/docker && docker build -t peted/warp-content .

docker-run: FORCE
	docker run -t -i -p 8080:8080 peted/warp-content

docker-push: FORCE
	docker push peted/warp-content

clean: FORCE
	@rm -rf ./misc/docker/target

prune: FORCE
	docker container prune -f
	docker image prune -f

FORCE:
