FROM golang:1.12.6 as builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
ARG CI_COMMIT_SHORT_SHA
ARG GITLAB_USERNAME
ARG GITLAB_PASSWORD
RUN echo "machine gitlab.skypicker.com " >> ~/.netrc &&\
  echo "  login $(curl -sS http://httpenv/v1/gitlab-user||echo $GITLAB_USERNAME)" >> ~/.netrc && \
  echo "  password $(curl -sS http://httpenv/v1/gitlab-password||echo $GITLAB_PASSWORD)" >> ~/.netrc && \
  make build &&\
  rm ~/.netrc

FROM alpine:3.9.4
RUN  apk add --no-cache --virtual=.run-deps ca-certificates &&\
  mkdir /app

WORKDIR /app
COPY --from=builder /app/main ./main
COPY --from=builder /app/.env.yaml .env.yaml
COPY --from=builder /app/.well-known .well-known/
COPY --from=builder /app/config config/
EXPOSE 8080
EXPOSE 7777

USER nobody
ENTRYPOINT ./main
