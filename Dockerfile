FROM golang:latest as builder
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main ./app
EXPOSE 8080

ENTRYPOINT ./app