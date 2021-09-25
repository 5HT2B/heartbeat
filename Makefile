NAME   := l1ving/heartbeat
TAG    := $(shell git rev-parse --short HEAD)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

heartbeat: clean update build

clean:
	rm -f heartbeat

update:
	go install github.com/valyala/quicktemplate/qtc
	go get -u github.com/valyala/fasthttp
	go get -u github.com/Ferluci/fast-realip
	go get -u golang.org/x/text/language
	go get -u golang.org/x/text/message

build:
	go generate
	go build -ldflags "-X main.gitCommitHash=$(TAG)" -o heartbeat .

docker-build:
	@docker build --build-arg COMMIT=${TAG} -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

docker-push:
	@docker push ${NAME}
