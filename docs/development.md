# Development

If you want to build the kernel from the source code, there are some requirements.

## Requirements

- ZeroMQ 4.2.3
- libzmq-dev
- Go 1.8 or above

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
