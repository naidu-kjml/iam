FROM golang:1.12.6-alpine

RUN apk add --no-cache git~=2.20 autoconf libtool curl gcc musl-dev zstd lz4

RUN curl -sSL https://github.com/ovh/venom/releases/download/v0.26.0/venom.linux-amd64 -o /usr/local/bin/venom
RUN chmod u+x /usr/local/bin/venom

COPY . /venom
WORKDIR /venom

CMD ["venom", "run", "--var-from-file", "variables.yml", "--parallel", "5", "--format=xml", "--output-dir", "."]
