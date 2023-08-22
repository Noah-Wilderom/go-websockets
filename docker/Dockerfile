FROM ubuntu:22.04

LABEL maintainer="Noah Wilderom"

WORKDIR /usr/bin/go-websockets

COPY ../ .

ENV TZ=Europe/Amsterdam

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update \
    && apt-get install -y gnupg curl ca-certificates libcap2-bin libpng-dev build-essential \
    && curl -sL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz | tar -C /usr/local -xz

ENV PATH="$PATH:/usr/local/go/bin"

EXPOSE 4001

RUN go build -o build/go-websockets main.go

RUN chmod +x build/go-websockets

RUN chmod +x docker/entrypoint.sh

ENTRYPOINT ["docker/entrypoint.sh"]