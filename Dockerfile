FROM golang:1.16.3  as build

WORKDIR /go/src/fileserver

COPY . .

RUN make static-install

FROM alpine:3.11

COPY --from=build /go/bin/* /usr/local/bin/

RUN apk update && apk --no-cache add ca-certificates && \
    update-ca-certificates
