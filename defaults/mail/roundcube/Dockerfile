FROM alpine:latest

RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.6/main" >> /etc/apk/repositories
RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.6/community" >> /etc/apk/repositories
RUN apk update
RUN apk --no-cache add \
      nginx \
      openssl \
      s6 \
      php7-fpm \
      php7-sqlite3 \
      php7-pdo \
      php7-zip \
      php7-pdo_sqlite \
      sqlite
RUN apk add roundcubemail==1.2.7-r0
RUN apk add roundcubemail-installer==1.2.7-r0
RUN pear install Mail_Mime
RUN pear install Mail_mimeDecode
RUN pear install Net_SMTP
RUN sed -ie 's^;daemonize = yes^daemonize = no^' /etc/php7/php-fpm.conf
RUN mkdir -p /etc/s6/.s6-svscan
RUN ln -s /bin/true /etc/s6/.s6-svscan/finish
RUN mkdir -p /etc/s6/nginx
RUN ln -s /bin/true /etc/s6/nginx/finish
RUN ln -s /usr/sbin/nginx /etc/s6/nginx/run
RUN mkdir -p /etc/s6/php-fpm
RUN ln -s /bin/true /etc/s6/php-fpm/finish
RUN ln -s /usr/sbin/php-fpm7 /etc/s6/php-fpm/run

COPY nginx.conf /etc/nginx/nginx.conf
COPY php.ini /etc/php7/php.ini
COPY config.inc.php /etc/roundcube/config.inc.php

EXPOSE 80
RUN mkdir /scripts
COPY start.sh /scripts/start.sh
RUN chmod +x /scripts/start.sh
RUN chown -R nobody:nobody /var/log/roundcube/
ENTRYPOINT ["/scripts/start.sh"]