NAME   := l1ving/heartbeat
TAG    := $(shell git log -1 --pretty=%h)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

heartbeat: clean
	go get -u github.com/valyala/fasthttp
	go get -u github.com/valyala/quicktemplate/qtc
	go get -u github.com/Ferluci/fast-realip
	go build -ldflags "-X main.gitCommitHash=$(TAG)"

clean:
	rm -f heartbeat

build:
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

push:
	@docker push ${NAME}

login:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}
