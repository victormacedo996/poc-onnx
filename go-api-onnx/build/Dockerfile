ARG GO_VERSION=1.23.2
ARG UBUNTU_VERSION=24.04
ARG ONNX_RUNTIME_VERSION=1.20.1

FROM ubuntu:${UBUNTU_VERSION} as builder
ARG GO_VERSION
ENV GO_VERSION=${GO_VERSION}

ARG ONNX_RUNTIME_VERSION
ENV ONNX_RUNTIME_VERSION=${ONNX_RUNTIME_VERSION} 

WORKDIR /go-api-onnx

COPY ./go-api-onnx /go-api-onnx
ENV PATH="${PATH}:/usr/local/go/bin"
ENV CGO_ENABLED=1

RUN apt update && \
    apt install \
        wget \
        gcc -y && \
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    wget https://github.com/microsoft/onnxruntime/releases/download/v${ONNX_RUNTIME_VERSION}/onnxruntime-linux-x64-${ONNX_RUNTIME_VERSION}.tgz && \
    tar -C /go-api-onnx -xzf onnxruntime-linux-x64-${ONNX_RUNTIME_VERSION}.tgz && \
    go build -o ./app ./cmd/api/main.go

FROM ubuntu:${UBUNTU_VERSION}
ARG GO_VERSION
ENV GO_VERSION=${GO_VERSION}

ARG ONNX_RUNTIME_VERSION
ENV ONNX_RUNTIME_VERSION=${ONNX_RUNTIME_VERSION} 

RUN apt update && apt upgrade -y && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /go-api-onnx
COPY --from=builder /go-api-onnx/app /go-api-onnx
COPY --from=builder /go-api-onnx/onnxruntime-linux-x64-${ONNX_RUNTIME_VERSION} /go-api-onnx/onnxruntime-linux-x64-${ONNX_RUNTIME_VERSION}
COPY ./models/logistic_regression_iris.onnx /models/

ENTRYPOINT [ "./app" ]