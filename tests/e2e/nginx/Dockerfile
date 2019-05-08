FROM nginx:1.16.0-alpine

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY ./mocks ./mocks

VOLUME ./mocks
