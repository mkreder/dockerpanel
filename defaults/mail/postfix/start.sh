#!/bin/sh

cp -rf /etc/postfix/amavisd.conf /etc/amavisd.conf
/usr/sbin/rsyslogd
postmap /etc/postfix/vmailbox
postmap /etc/postfix/virtual
postmap /etc/postfix/transport

newaliases
/usr/sbin/postfix -c /etc/postfix start
/usr/sbin/amavisd

# Ya que postfix forkea y solo corre en background tenemos que hacer algo para que el contenedor siga corriendo
while true; do
    sleep 1000
done