#!/bin/bash

set -x

if [ -z "${IMAGE_REGISTRY}" ]; then
	echo "IMAGE_REGISTRY is not set"
	exit 1
fi

TARBALL=${1:-}
tar -zxvf $TARBALL

CMD="./crane"

$CMD push images/appscode-cluster-ui-0.9.7.tar $IMAGE_REGISTRY/appscode/cluster-ui:0.9.7
