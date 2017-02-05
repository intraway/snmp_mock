#!/bin/sh

SNMP_CONF_FILE=/app/snmpd.conf
MOCK_CONF_FILE=/app/config.yaml

if [ ! -z "$SNMP_PORT" ];then
    echo "Replacing snmp port with $SNMP_PORT"
    sed -i "s/agentAddress udp:\([0-9][0-9]*\)/agentAddress udp:$SNMP_PORT/" $SNMP_CONF_FILE
fi

if [ ! -z "$SNMP_COMMUNITY" ];then
    echo "Replacing snmp community with $SNMP_COMMUNITY"
    sed -i "s/rocommunity .*/rocommunity $SNMP_COMMUNITY/" $SNMP_CONF_FILE
fi

if [ ! -z "$BASE_OID" ];then
    echo "Replacing base_oid with $BASE_OID"
    sed -i "s/base_oid:.*/base_oid: $BASE_OID/" $MOCK_CONF_FILE
fi

if [ ! -z "$APP_PORT" ];then
    echo "Replacing app_port with $APP_PORT"
    sed -i "s/app_port:.*/app_port: $APP_PORT/" $MOCK_CONF_FILE
fi

/usr/sbin/snmpd -Lsd -C -c $SNMP_CONF_FILE
/app/snmp_mock `find /app/oids -type f -name "*.csv"`
