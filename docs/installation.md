# Installation

## Dockerized Kernel

### Requirements

- Jupyter Notebook v4.4.0 or above (We do not test other versions now.)
- Kubernetes v1.9 or above
- Docker

### Update config.yaml

Update `./hack/config.yaml`:

```
kubeconfig: /var/run/kubernetes/admin.kubeconfig
s2i:
  provider: img
  username: <input-your-username>
  password: <input-your-pwd>
```

### Build the Docker Image

Run the command:

```
docker build -t caicloud/ciao .
```

### Install Dockerized kernel

Run the command:

```
./hack/install-dockerized.sh
```

## Native

### Requirements

- Jupyter Notebook v4.4.0 or above (We do not test other versions now.)
- Kubernetes v1.9 or above
- [Kubeflow](https://www.kubeflow.org/) v0.2 or above
- [S2I](https://github.com/openshift/source-to-image) v1.1.10 or above
- [img](https://github.com/genuinetools/img)
- Docker

### Get the kernel

We do not release any version of the kernel, so please build it from the source, please see the [Development Guide](./development.md).

### Setup Kubernetes and Kubeflow

Please see [Getting Started with Kubeflow](https://www.kubeflow.org/docs/started/getting-started/). Currently, we only support one-node cluster since we do not push the image to the Docker Registry.

### Install the Kernel

Run the script `hack/install.sh`, then the specification of the kernel will be installed to `${HOME}/.local/share/jupyter/kernels/kubeflow`, then Jupyter will know the information about Ciao.

Then we need to create a configuration file `$HOME/.ciao/config.yaml`:

```yaml
kubeconfig: {path to your kubeconfig}
s2i:
  provider: {img or s2i}
  registry: {registry to be used to push images, optional}
  username: {username}
  password: {password}
```

There are two options about tools to convert the source code in Jupyter Notebook to Docker image:

- [img](https://github.com/genuinetools/img) (Recommended), which is a daemon-free tool to build and push Docker images.
- [s2i](https://github.com/openshift/source-to-image), which is a source to image tool.

### Install Image

For better performance, we recommend pulling the builder images from Docker Registry ahead of time. There are two builder images for different ML frameworks:

- `gaocegege/tensorflow-s2i:1.10.1-py3`
- `gaocegege/pytorch-s2i:v0.2`
- `tensorflow/tensorflow:1.10.1-py3`
- `pytorch/pytorch:v0.2`

Or the time of the first run will be extremely long (which depends on your network).

### Run the Kernel

First, we need to set the environment variable `KUBECONFIG` to tell the kernel where to find the kubeconfig:

```bash
export KUBECONFIG={path to your kubeconfig}
```

Then run Jupyter Notebook or Lab, choose Kubeflow kernel.
