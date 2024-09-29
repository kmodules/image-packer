#!/bin/bash

set -x

if [ -z "${IMAGE_REGISTRY}" ]; then
	echo "IMAGE_REGISTRY is not set"
	exit 1
fi

OS=$(uname -o)
if [ "${OS}" = "GNU/Linix" ]; then
  OS=Linux
fi
ARCH=$(uname -m)
if [ "${ARCH}" = "aarch64" ]; then
  ARCH=arm64
fi
curl -sL "https://github.com/google/go-containerregistry/releases/latest/download/go-containerregistry_${OS}_${ARCH}.tar.gz" > /tmp/go-containerregistry.tar.gz
tar -zxvf /tmp/go-containerregistry.tar.gz -C /tmp/
mv /tmp/crane .

CMD="./crane"

$CMD cp ghcr.io/appscode/cluster-ui:0.9.7 $IMAGE_REGISTRY/appscode/cluster-ui:0.9.7
