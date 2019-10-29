REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.revision=$(REVISION)"

build:
	# protoc -I=./storage --go_out=./storage ./storage/data.proto
	go build -ldflags $(LDFLAGS)

run : build
	go run .

fmt :
	go fmt ./...
	git diff --stat

test : build
	go test ./...

clean:
	rm -rf .bogo
