FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
ARG CI_COMMIT_SHORT_SHA
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main ./app
COPY --from=builder /app/.env.yaml .env.yaml
EXPOSE 8080

ENTRYPOINT ./app