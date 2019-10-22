build:
	protoc -I=./storage --go_out=./storage ./storage/data.proto
	go build

run : build
	go run .

clean:
	rm -rf .bogo
