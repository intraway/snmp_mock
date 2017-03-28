FROM alpine:3.4

ARG BUILD_DATE
ARG BUILD_VCS_REF
ARG BUILD_VERSION


LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-url="https://github.com/intraway/snmp_mock.git" \
      org.label-schema.vcs-ref=$BUILD_VCS_REF \
      org.label-schema.version=$BUILD_VERSION \
      com.microscaling.license=GPL-3.0



ADD docker /app
#ADD snmp_mock /app/snmp_mock
ADD sample_oids /app/oids

RUN apk add --no-cache net-snmp

ENV GOROOT=/usr/lib/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin

WORKDIR /gopath/src/app
ADD . /gopath/src/app

RUN apk add --no-cache git go g++ && \
  go get && \
  go build && \
  cp app /app/snmp_mock && \
  apk del git go g++ && \
  rm -rf /gopath

WORKDIR /app
ENTRYPOINT ["/app/init.sh"]
