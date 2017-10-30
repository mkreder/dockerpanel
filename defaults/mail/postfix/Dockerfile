FROM alpine:latest
RUN apk update
RUN apk add postfix rsyslog spamassassin amavisd-new perl-socket6 perl-io-socket-inet6
RUN sa-update
RUN mkdir /scripts
COPY start.sh /scripts/start.sh
RUN chmod +x /scripts/start.sh
VOLUME	[ "/var/spool/postfix", "/var/mail", "/etc/postfix" ]
EXPOSE 25
ENTRYPOINT ["/scripts/start.sh"]