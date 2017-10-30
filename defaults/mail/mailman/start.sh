#!/bin/sh


if ! [ -d /usr/lib/mailman/Mailman ]; then
  tar xvf usrlib.tar -C /
fi

if ! [ -d /var/lib/mailman/Mailman ]; then
  tar xvf varlib.tar -C /
fi

/etc/init.d/rsyslog start
chmod +x /etc/mailman/add_lists.sh
/etc/mailman/add_lists.sh
rm -rf /etc/mailman/add_lists.sh
/usr/lib/mailman/bin/genaliases  > /var/lib/mailman/data/virtual-mailman
/usr/lib/mailman/bin/mailmanctl start


while true ; do
  sleep 1000
done