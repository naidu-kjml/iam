FROM golang:1.12.6-alpine

RUN apk add --no-cache git~=2.20 autoconf libtool curl gcc musl-dev zstd lz4

RUN curl -sSL https://github.com/ovh/venom/releases/download/v0.25.0/venom.linux-amd64 -o /usr/local/bin/venom
RUN chmod u+x /usr/local/bin/venom

RUN curl -sSL https://github.com/fullstorydev/grpcurl/releases/download/v1.2.1/grpcurl_1.2.1_linux_x86_64.tar.gz | tar -xz -C /usr/local/bin
RUN chmod u+x /usr/local/bin/grpcurl

COPY . /venom
WORKDIR /venom

CMD ["venom", "run", "--var-from-file", "variables.yml", "--parallel", "5", "--format=xml", "--output-dir", "."]
