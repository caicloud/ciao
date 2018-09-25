#!/usr/bin/env bash

ROOT=$(dirname ${BASH_SOURCE})/..

cd ${ROOT}
rm -r ${HOME}/.local/share/jupyter/kernels/kubeflow
mkdir -p ${HOME}/.local/share/jupyter/kernels/kubeflow
cp -r ./hack/kernel.json ${HOME}/.local/share/jupyter/kernels/kubeflow/kernel.json
cd - > /dev/null
