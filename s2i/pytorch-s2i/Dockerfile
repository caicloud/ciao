FROM pytorch/pytorch:v0.2

LABEL maintainer="Kubeflow"

LABEL io.k8s.description="Kubeflow image for PyTorch" \
    io.k8s.display-name="PyTorch" \
    # this label tells s2i where to find its mandatory scripts
    # (run, assemble, save-artifacts)
    io.openshift.s2i.scripts-url="image:///usr/libexec/s2i"

# Copy the S2I scripts to /usr/libexec/s2i since we set the label that way
COPY ./s2i/bin /usr/libexec/s2i

RUN mkdir /app
