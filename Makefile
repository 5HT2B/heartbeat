heartbeat: clean
	go get -u github.com/valyala/fasthttp
	go get -u github.com/valyala/quicktemplate/qtc
	go build

clean:
	rm -f heartbeat
