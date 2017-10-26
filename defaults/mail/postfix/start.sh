#!/bin/sh

postmap /etc/postfix/vmailbox
/usr/sbin/rsyslogd
/usr/sbin/postfix -c /etc/postfix start

# Ya que postfix forkea y solo corre en background tenemos que hacer algo para que el contenedor siga corriendo
while true; do
    sleep 1000
done