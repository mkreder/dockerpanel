#!/bin/sh

if [ ! -d "/run/mysqld" ]; then
    mkdir -p /run/mysqld
    chown -R mysql:mysql /run/mysqld
fi

if [ -d /var/lib/mysql/mysql ]; then
    echo "MySQL ya configurado"
else
   chown -R mysql:mysql /var/lib/mysql
   mysql_install_db --user=mysql
   if [ -f userdata.sql ]; then
       /usr/bin/mysqld --user=mysql --bootstrap --verbose=0 < userdata.sql
   fi
fi
/usr/bin/mysqld --user=mysql --console
