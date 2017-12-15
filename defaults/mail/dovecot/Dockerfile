FROM alpine:latest
RUN apk update
RUN apk add dovecot rsyslog dovecot-pigeonhole-plugin bash
RUN mkdir /scripts
COPY start.sh /scripts/start.sh
RUN chmod +x /scripts/start.sh
RUN mkdir /var/mail
RUN addgroup -g 104 vmail
RUN adduser -h /var/mail -G vmail -u 101 -H vmail -D
VOLUME	[ "/var/mail", "/etc/dovecot" ]
EXPOSE 110 143 12345 12340
ENTRYPOINT ["/scripts/start.sh"]