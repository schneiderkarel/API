FROM golang:1.16.0-alpine

RUN apk update \
  && apk --update add git openssl openssh-client bash curl \
  && go get github.com/githubnemo/CompileDaemon

COPY cmd/docker/local/docker-entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /app
COPY . .

EXPOSE 8080

ENTRYPOINT ["entrypoint.sh"]
