build:
	protoc -I=./storage --go_out=./storage ./storage/data.proto
	go build

clean:
	rm bogodb catalog.db
