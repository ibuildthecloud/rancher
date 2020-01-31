#!/bin/bash
set -x -e
#go build -o package/agent ./pkg/agent/
go build -o package/rancher
#docker build -t ibuildthecloud/rancher-agent:dev -f package/Dockerfile.agent package/
docker build -t ibuildthecloud/rancher:dev -f package/Dockerfile package
#docker push ibuildthecloud/rancher-agent:dev
docker push ibuildthecloud/rancher:dev
