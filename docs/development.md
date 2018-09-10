# Development

If you want to build the kernel from the source code, there are some requirements.

## Requirements

- ZeroMQ 4.2.3
- libzmq-dev
- Go 1.8 or above

## Set KUBECONFIG

We need to set KUBECONFIG manually here [main.go](https://github.com/caicloud/ciao/blob/4ce8580110c3ec825cefc2e05abfbf6150d06e02/cmd/kubeflow-kernel/main.go#L30):

```go
    // This is a trick! Jupyter does not use the env var to create kernel.
	kubeconfig := "/var/run/kubernetes/admin.kubeconfig"
```

We need to support configuration just like IPython, but now we have to do this dirty work.

## Build

Make sure that `$GOPATH/bin` is in the `$PATH`:

```bash
export PATH=$PATH:"${GOPATH}/bin"
```

Then install the kernel by this command:

```bash
go install github.com/caicloud/ciao/cmd/kubeflow-kernel
```

To verify the binary is installed successfully, when you run `kubeflow-kernel`, you should get a log:

```text
Need a command line argument specifying the connection file.
```
