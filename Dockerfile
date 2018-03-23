FROM alpine
MAINTAINER Christopher Hein <me@christopherhein.com>

RUN apk --no-cache add openssl musl-dev ca-certificates
ADD node-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/node-operator"]
