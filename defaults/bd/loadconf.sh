#!/bin/sh

while ! [ -S /run/mysqld/mysqld.sock ]; do
    echo "MySQL todavia no esta listo, esperando"
    sleep 1
done
mysql < /conf/userdata.sql
