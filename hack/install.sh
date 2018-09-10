#!/usr/bin/env bash

ROOT=$(dirname ${BASH_SOURCE})/..

cd ${ROOT}
rm -r ${HOME}/.local/share/jupyter/kernels/kubeflow
cp -r ./artifacts ${HOME}/.local/share/jupyter/kernels/kubeflow
cd - > /dev/null
