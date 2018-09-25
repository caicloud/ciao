package img

const (
	tensorflowTemplate = `FROM tensorflow/tensorflow:1.10.1-py3
LABEL maintainer="Kubeflow"
COPY ./code.py /
ENTRYPOINT ["python", "/code.py"]`

	pytorchTemplate = `FROM pytorch/pytorch:v0.2
LABEL maintainer="Kubeflow"
COPY ./code.py /
ENTRYPOINT ["python", "/code.py"]`
)
