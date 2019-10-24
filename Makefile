build:
	protoc -I=./storage --go_out=./storage ./storage/data.proto
	go build

run : build
	go run .

fmt :
	go fmt ./...
	git diff --stat

clean:
	rm -rf .bogo
