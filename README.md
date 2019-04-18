# Ciao - Kernel for Kubeflow in Jupyter Notebook

[![Go Report Card](https://goreportcard.com/badge/github.com/caicloud/ciao)](https://goreportcard.com/report/github.com/caicloud/ciao)
[![Build Status](https://travis-ci.org/caicloud/ciao.svg?branch=master)](https://travis-ci.org/caicloud/ciao)
[![Coverage Status](https://coveralls.io/repos/github/caicloud/ciao/badge.svg?branch=master)](https://coveralls.io/github/caicloud/ciao?branch=master)

Ciao is still in early development -- it is not feature-complete or production-ready. Please try our experimental kernel and give us your feedback.

- Report bugs to [Ciao Issues](https://github.com/caicloud/ciao/issues)
- Chat at the [Kubeflow Slack Channel](https://kubeflow.slack.com/messages/CJ01RLD7Y) (Please @gaocegege to let us know)

## Overview

Ciao is a jupyter kernel for Kubeflow. The name of the project `Ciao` comes from Italian:

> The word "ciao" (/ˈtʃaʊ/; Italian pronunciation: [ˈtʃaːo]) is an informal salutation in the Italian language that is used for both "hello" and "goodbye".

Ciao's goal is to simplify the machine learning workflow using Kubeflow. Currently, users could create a **distributed** model training job from Jupyter Notebook and get the logs of all replicas (parameter servers and workers) in the output.

## Demo

Please see the [Demo Show](./docs/demo.md).

![Ciao and SOS integration](./docs/images/ciao-sos-integration.gif)

### Usage

There are some magic commands supported by Ciao:

```
%framework=tensorflow
%ps={number}
%worker={number}
```

or

```
%framework=pytorch
%master={number}
%worker={number}
```

### Examples

- [Distributed Training using TensorFlow](./docs/examples/tensorflow/example.ipynb)
- [Distributed Training using PyTorch](./docs/examples/tensorflow/example.ipynb)

## Installation

Please see the [Installation Guide](./docs/installation.md).

## Design Document

Please see the [Design Document](./docs/design.md) to know the architecture of Ciao.

## Acknowledgments

- Thank [kubeflow/kubeflow](https://github.com/kubeflow/kubeflow) for the awesome operators which supports TensorFlow/PyTorch and many other ML frameworks on Kubernetes.
- Thank [gopherdata/gophernotes](https://github.com/gopherdata/gophernotes) for the reference implementation of Jupyter Kernel in Golang.
