FROM tensorflow/tensorflow:1.10.1-py3

LABEL maintainer="Kubeflow"

LABEL io.k8s.description="Kubeflow image for TensorFlow" \
    io.k8s.display-name="TensorFlow" \
    # this label tells s2i where to find its mandatory scripts
    # (run, assemble, save-artifacts)
    io.openshift.s2i.scripts-url="image:///usr/libexec/s2i"

# Copy the S2I scripts to /usr/libexec/s2i since we set the label that way
COPY ./s2i/bin /usr/libexec/s2i

RUN mkdir /app
