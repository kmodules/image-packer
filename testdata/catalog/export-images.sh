#!/bin/bash
set -x

mkdir -p images

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
mv /tmp/crane images

CMD="./images/crane"

$CMD pull ghcr.io/appscode/cluster-ui:0.9.7 images/appscode-cluster-ui-0.9.7.tar

tar -czvf images.tar.gz images
