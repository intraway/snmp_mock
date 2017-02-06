FROM alpine

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
