#!/bin/sh
if [ ! -f /usr/share/webapps/roundcube/config/config.inc.php ]; then
  cp /etc/roundcube/* /usr/share/webapps/roundcube/config
fi

chown -R nobody:nobody /usr/share/webapps/roundcube/config \
  && exec s6-svscan /etc/s6

while true; do
 sleep 1000
 done