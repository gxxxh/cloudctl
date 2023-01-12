# Build the manager binary
FROM golang:1.19 as builder
ENV GOPROXY "https://goproxy.cn,direct"

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY src/ src/
COPY config/ config/
ENV CLOUDCTL_CONFIG_PATH /workspace/config/crd_configs/openstack

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN go build -o /usr/local/go/bin/cloudctl

ENTRYPOINT ["cloudctl"]
#CMD sleep 10000000
