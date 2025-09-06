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
	go get -u github.com/redis/go-redis/v9
	go get -u github.com/joho/godotenv
	go get -u github.com/nitishm/go-rejson/v4
	go get -u github.com/valyala/fasthttp
	go get -u github.com/valyala/quicktemplate
	go get -u golang.org/x/text/language
	go get -u golang.org/x/text/message

docker-build:
	@docker build --build-arg COMMIT=${TAG} -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

docker-push:
	@docker push ${NAME}

lint:
	@echo "Running pre-commit hooks..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit run --all-files; \
	else \
		echo "pre-commit not found. Install with: pip install pre-commit"; \
		exit 1; \
	fi
