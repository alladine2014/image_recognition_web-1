#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
if [ "X$1" != "X" ]; then
    RUNTIME_ROOT=$1
else
    RUNTIME_ROOT=${CURDIR}
fi

export GODEBUG=netdns=cgo
export ENV_PSM="image.recognition.web"

exec ${CURDIR}/bin/image-recognition -addr 0.0.0.0 -port 8080 -conf ${CURDIR}/conf
