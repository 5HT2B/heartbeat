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
	go get -u github.com/ferluci/fast-realip
	go get -u github.com/go-redis/redis/v8
	go get -u github.com/joho/godotenv
	go get -u github.com/nitishm/go-rejson/v4
	go get -u github.com/valyala/fasthttp
	go get -u github.com/valyala/quicktemplate
	go get -u golang.org/x/text/language
	go get -u golang.org/x/text/message
	go get -u github.com/fasthttp/websocket

docker-build:
	@docker build --build-arg COMMIT=${TAG} -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

docker-push:
	@docker push ${NAME}
