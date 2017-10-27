FROM debian
# Se usa debian porque alpine no cuenta con mailman
ENV  DEBIAN_FRONTEND noninteractive
RUN apt-get update
RUN apt-get --asume-yes true install mailman apache2 rsyslog
COPY apache2.conf  /etc/apache2/sites-available/mailman.conf
RUN mkdir /scripts
COPY start.sh /scripts
RUN chmod +x /scripts/start.sh
RUN mkdir /var/www/lists
RUN a2ensite mailman.conf
RUN newlist --emailhost=ejemplo.com -a  mailman mailman@ejemplo.com mailman
RUN mkdir /var/run/mailman/
RUN chown -R list:list /var/run/mailman/
WORKDIR /root
RUN tar cvf usrlib.tar /usr/lib/mailman
RUN tar cvf varlib.tar /var/lib/mailman
ENTRYPOINT ["/scripts/start.sh"]
EXPOSE 80