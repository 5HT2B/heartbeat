NAME   := l1ving/heartbeat
TAG    := $(shell git rev-parse --short HEAD)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

heartbeat: clean deps build

clean:
	rm -f heartbeat

generate:
	go generate

build: generate
	go build -ldflags "-X main.gitCommitHash=$(TAG)" -o heartbeat .

deps:
	go install github.com/valyala/quicktemplate/qtc
	go get -u ./...

docker-build:
	@docker build --build-arg COMMIT=${TAG} -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

docker-push:
	@docker push ${NAME}
