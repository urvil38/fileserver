LDFLAGS=v1.0.0

linux-build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -X main.version=${LDFLAGS}" .

macos-build:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -X main.version=${LDFLAGS}" .

windows-build:
	GOOS=windows GOARCH=amd64 go build -ldflags="-w -X main.version=${LDFLAGS}" .
