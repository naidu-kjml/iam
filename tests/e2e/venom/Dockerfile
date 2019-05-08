FROM golang:1.12.4-alpine

RUN apk add --no-cache git~=2.20 autoconf libtool curl gcc musl-dev

RUN curl -sSL https://github.com/ovh/venom/releases/download/v0.25.0/venom.linux-amd64 -o /usr/local/bin/venom
RUN chmod u+x /usr/local/bin/venom

COPY . /venom
WORKDIR /venom

CMD ["venom", "run", "--var-from-file", "variables.yml", "--parallel", "5", "--format=xml", "--output-dir", "."]
