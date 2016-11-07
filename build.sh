#!/bin/bash

rm -rf bin && mkdir bin && cd bin
echo "clean and make bin"

echo "build ruleng (linux x64) ..."
GOOS=linux GOARCH=amd64 go build ../. && tar -zcvf ruleng_linux_amd64.tar.gz ruleng && rm -rf ruleng

echo "build all arch Successfully"
