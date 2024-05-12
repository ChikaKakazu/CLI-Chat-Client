FROM golang:1.22.3-alpine3.19

RUN apk update && apk add git curl unzip

ENV GOPATH /go

RUN mkdir /go/app
COPY . /go/app

WORKDIR /go/app

# バージョン取得
# $(curl -s "https://api.github.com/repos/protocolbuffers/protobuf/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')
ENV PROTOC_VERSION 26.1
ENV PROTOC_ZIP protoc-${PROTOC_VERSION}-linux-x86_64.zip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP} 
RUN unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc \
    && unzip -o ${PROTOC_ZIP} -d /usr/local 'include/*' \
    && rm -f ${PROTOC_ZIP}

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go get -u google.golang.org/grpc@latest
