FROM alpine:latest

MAINTAINER "River Yang" "comicme_yanghe@nanoframework.org"

# RUN apk update && apk --no-cache add ca-certificates wget

# RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub

# RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.23-r3/glibc-2.23-r3.apk

# RUN apk add glibc-2.23-r3.apk

RUN apk update && apk add socat

ADD /bin/ruleng_linux_amd64.tar.gz /

COPY entrypoint.sh /

EXPOSE 8080

CMD ["/entrypoint.sh"]
