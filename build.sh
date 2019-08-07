#!/bin/bash -e
RUN_NAME="image.recognition.web"
GIT_SHA=`git rev-parse --short HEAD || echo "NotGitVersion"`

WHEN=`date '+%Y-%m-%d_%H:%M:%S'`

go build -v -ldflags "-s -X main.GitSHA=${GIT_SHA} -X main.BuildTime=${WHEN}" -o output/bin/image-recognition 

mkdir -p output/bin output/conf output/log
cp script/bootstrap.sh output
cp script/settings.py output
cp -rf data output/
chmod +x output/bootstrap.sh
chmod +x output/settings.py
cp -rf conf/* output/conf/
