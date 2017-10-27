#!/bin/sh


if ! [ -d /usr/lib/mailman/Mailman ]; then
  tar xvf usrlib.tar -C /
fi

if ! [ -d /var/lib/mailman/Mailman ]; then
  tar xvf varlib.tar -C /
fi

/etc/init.d/rsyslog start
/usr/lib/mailman/bin/genaliases  > /var/lib/mailman/data/virtual-mailman
/usr/lib/mailman/bin/mailmanctl start
chmod +x /etc/mailman/run.sh
/etc/mailman/run.sh
rm -rf /etc/mailman/run.sh


while true ; do
  sleep 1000
done