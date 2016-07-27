#!/bin/bash

PD_ROOT=$(dirname "${BASH_SOURCE[0]}")
pushd `dirname $0` > /dev/null
PD_ROOT=$(cd $PD_ROOT && cd .. && pwd)
popd > /dev/null

rm -rf ${PD_ROOT}/vendor && ln -s ${PD_ROOT}/_vendor/vendor ${PD_ROOT}/vendor

CGO_ENABLED=0 go build -v -a -installsuffix cgo -o ${PD_ROOT}/kubernetes/docker-q8sad/pd-server ${PD_ROOT}/kubernetes/pd-server/main.go

rm -rf ${PD_ROOT}/vendor
