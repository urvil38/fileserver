LDFLAGS=1.0.0
VERSION=v1.0
LIST=linux darwin
build:
	@for os in $(LIST); do \
		CGO_ENABLED=0 GOOS=$$os GOARCH=amd64 go build -ldflags="-s -X main.version=${LDFLAGS}" . ; \
		mkdir -p ~/Documents/fileserver-bin/${VERSION}/$$os ; \
		mv fileserver  ~/Documents/fileserver-bin/${VERSION}/$$os ; done \

	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -X main.version=${LDFLAGS}" .
	mkdir -p ~/Documents/fileserver-bin/${VERSION}/windows ;
	mv fileserver.exe  ~/Documents/fileserver-bin/${VERSION}/windows

upload:
	cd ~/go/src/fileserver/upload && go run main.go -v ${VERSION}
