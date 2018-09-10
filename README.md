# Ciao - Kernel for Kubelow in Jupyter Notebook

Ciao is still in early development -- it is not feature-complete or production-ready. Please try our experimental kernel and give us your feedback.

- Report bugs to [Ciao Issues](https://github.com/caicloud/ciao/issues)
- Chat at the [Kubeflow Slack Channel](https://kubeflow.slack.com/messages/C7REE0EHK/) (Please @gaocegege to let us know)

## Overview

Ciao is a jupyter kernel for Kubeflow. The name of the project `Ciao` comes from Italian:

> The word "ciao" (/ˈtʃaʊ/; Italian pronunciation: [ˈtʃaːo]) is an informal salutation in the Italian language that is used for both "hello" and "goodbye".

Ciao's goal is to simplify the machine learning workflow using Kubeflow. Currently, users could create a **distributed** model training job from Jupyter Notebook and get the logs of all replicas (parameter servers and workers) in the output.

## Demo

Please see the [Demo Show](./docs/demo.md).

### Usage

There are some magic commands supported by Ciao:

- `%kubeflow framework={framework that you want to use}` This command defines which framework will be used. Currently we only support tensorflow.
- `%kubeflow ps={number of ps}` This command defines how many parameter servers you want to create.
- `%kubeflow worker={number of workers}` This command defines how many workers you want to create.

### Examples

- [Distributed Training with Ciao](./docs/examples/example.ipynb)

## Installation

Please see the [Installation Guide](./docs/installation.md).

## Design Document

Please see the [Design Document](./docs/design.md) to know the architecture of Ciao.

## Acknowledgments

- Thank [kubeflow/kubeflow](https://github.com/kubeflow/kubeflow) for the awesome operators which supports TensorFlow/PyTorch and many other ML frameworks on Kubernetes.
- Thank [gopherdata/gophernotes](https://github.com/gopherdata/gophernotes) for the reference implementation of Jupyter Kernel in Golang.
- Thank [openshift/source-to-image](https://github.com/openshift/source-to-image) for the tool to convert source code to Docker image directly.
