FROM alpine:latest
RUN apk update
RUN apk add mysql mysql-client
RUN mkdir /scripts
VOLUME ["/conf","/var/lib/mysql"]
COPY startup.sh /scripts/startup.sh
COPY loadconf.sh /scripts/loadconf.sh
RUN chmod +x /scripts/*.sh
ENTRYPOINT ["/scripts/startup.sh"]
EXPOSE 3306
