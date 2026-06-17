.PHONY: bin clean lint test tidy run cleandb

DEFAULT_GOAL := bin

bin:
	go build -o build/nsf main.go

clean:
	rm -rf build

lint:
	golangci-lint run

test:
	go test -v ./...

tidy:
	go mod tidy

run:
	./build/nsf -c config.yaml

cleandb:
	rm /tmp/nsf.db

cleantestdb:
	rm /tmp/nsf_db_test.db
	rm /tmp/nsf_db_context_test.db