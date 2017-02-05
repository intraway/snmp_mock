FROM alpine

ADD docker /app
ADD snmp_mock /app/snmp_mock
ADD sample_oids /app/oids

RUN apk add --no-cache net-snmp

ENTRYPOINT ["/app/init.sh"]
