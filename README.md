# Ciao - Kernel for Kubelow in Jupyter Notebook

This is an experimental project.

## Overview

Ciao is a jupyter kernel for Kubeflow. The name of the project `Ciao` comes from Italian:

> The word "ciao" (/ˈtʃaʊ/; Italian pronunciation: [ˈtʃaːo]) is an informal salutation in the Italian language that is used for both "hello" and "goodbye".

Ciao's goal is to simplify the machine learning workflow using Kubeflow. Currently, users could create a **distributed** model training job from Jupyter Notebook and get the logs of all replicas (parameter servers and workers) in the output.

## Demo

Please see the [demo page](./docs/demo.md)

### Usage

There are some magic commands supported by Ciao:

- `%kubeflow framework={framework that you want to use}` This command defines which framework will be used. Currently we only support tensorflow.
- `%kubeflow ps={number of ps}` This command defines how many parameter servers you want to create.
- `%kubeflow worker={number of workers}` This command defines how many workers you want to create.

### Examples

- [Distributed Training With Ciao](./docs/examples/example.ipynb)

### Installation

Please see the [installation guide](./docs/installation.md).
