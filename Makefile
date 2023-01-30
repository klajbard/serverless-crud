.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/upload ./handlers/upload/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/findall ./handlers/findall/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/find ./handlers/find/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/remove ./handlers/remove/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update ./handlers/update/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
