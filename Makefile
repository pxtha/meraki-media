PROJECT_NAME=meradia
BUILD_VERSION=0.0.1

DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin main.go

compose_dev: docker
	cd deploy && BUILD_VERSION=$(BUILD_VERSION) docker-compose up --build --force-recreate -d

compose_prd: docker
	cd deploy && docker tag $(PROJECT_NAME):$(BUILD_VERSION) gcr.io/joon-app/$(PROJECT_NAME):latest && docker push gcr.io/joon-app/$(PROJECT_NAME):latest

docker_prebuild: build
	mkdir -p deploy/conf
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deploy/$(PROJECT_NAME).bin; \
	cp -R conf deploy/;

docker_build:
	cd deploy; \
	docker build --rm -t $(DOCKER_IMAGE) .;

docker_postbuild:
	cd deploy; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\
	rm -rf conf 2> /dev/null;
docker: docker_prebuild docker_build docker_postbuild

mock:
	go generate -x -run="mockgen" ./...


