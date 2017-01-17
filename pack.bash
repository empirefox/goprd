#!/bin/sh
set -e

docker run --rm -it \
	-v "$GOPATH":/gopath \
	-v "$(pwd)":/app \
	-e "GOPATH=/gopath" \
	-w /app \
	golang:1.7.4-alpine \
	sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s -w" -o main'

rm -f main-*.tar.gz
sha256=`tar -zcvf - ./main ./xps-config.json ./xps-files/xps-prod.tar.gz | tee main.tar.gz | sha256sum | awk '{ print $1 }'`
mv main.tar.gz main-${sha256}.tar.gz