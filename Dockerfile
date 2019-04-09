FROM golang:1.12.3 as builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
ARG CI_COMMIT_SHORT_SHA
RUN make build

FROM alpine:3.9.2
RUN  apk add --no-cache --virtual=.run-deps ca-certificates &&\
  mkdir /dist/

WORKDIR /dist/
COPY --from=builder /app/main ./app
COPY --from=builder /app/.env.yaml .env.yaml
EXPOSE 8080

USER nobody
ENTRYPOINT ./app
