build-hash=$(shell git log --pretty=%h | head -n 1)

heartbeat: clean
	go get -u github.com/valyala/fasthttp
	go get -u github.com/valyala/quicktemplate/qtc
	go build -ldflags "-X main.gitCommitHash=$(build-hash)"

clean:
	rm -f heartbeat
