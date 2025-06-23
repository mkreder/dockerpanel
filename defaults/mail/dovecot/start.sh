#!/bin/bash

cd /etc/dovecot/auth/
for domain in *; do
    if ! [ -d /var/mail/vhosts/${domain} ]; then
      mkdir -p /var/mail/vhosts/${domain}
    fi
    cp -rf $domain/passwd  /var/mail/vhosts/$domain/
    cp -rf $domain/shadow  /var/mail/vhosts/$domain/
    for account in $domain/*; do
        if [ "$account" != "${domain}/passwd" ] && [ "$account" != "${domain}/shadow" ]; then
            if ! [ -d /var/mail/vhosts/$account ]; then
              mkdir -p /var/mail/vhosts/$account
            fi
            cp -rf $account/autoreply.sieve  /var/mail/vhosts/$account/
        fi
    done
done
chown -R 101:104 /var/mail/vhosts/

/usr/sbin/rsyslogd
/usr/sbin/dovecot -F