#!/usr/bin/env bash

ROOT=$(dirname ${BASH_SOURCE})/..

cd ${ROOT}
cp -r ./artifacts ${HOME}/.local/share/jupyter/kernels/kubeflow
cd - > /dev/null
