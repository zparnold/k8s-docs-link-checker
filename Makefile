build:
	dep ensure
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/checker src/main.go