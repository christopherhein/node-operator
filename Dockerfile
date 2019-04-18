FROM alpine
MAINTAINER Chris Hein <me@chrishein.com>

RUN apk --no-cache add openssl musl-dev ca-certificates libc6-compat
ADD node-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/node-operator"]
