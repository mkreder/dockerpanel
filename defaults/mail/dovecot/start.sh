#!/bin/sh

cd /etc/dovecot/auth/
for dominio in *; do
    if ! [ -d /var/mail/vhosts/${dominio} ]; then
      mkdir -p /var/mail/vhosts/${dominio}
    fi
    cp -rf $dominio/passwd  /var/mail/vhosts/$dominio/
    cp -rf $dominio/shadow  /var/mail/vhosts/$dominio/
    for cuenta in $dominio/*; do
        if [ "$cuenta" != "${dominio}/passwd" ] && [ "$cuenta" != "${dominio}/shadow" ]; then
            if ! [ -d /var/mail/vhosts/$dominio/$cuenta ]; then
              mkdir -p /var/mail/vhosts/$dominio/$cuenta
            fi
            cp -rf $dominio/$cuenta/autoreply.sieve  /var/mail/vhosts/$dominio/$cuenta/
        fi
    done
done
chown -R 101:104 /var/mail/vhosts/

/usr/sbin/rsyslogd
/usr/sbin/dovecot -F