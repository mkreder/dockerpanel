FROM alpine
RUN apk update
RUN apk add vsftpd
RUN echo "anon_root=/var/lib/ftp/" >> /etc/vsftpd/vsftpd.conf
RUN sed -i "/write_enable/s/#write/write/" /etc/vsftpd/vsftpd.conf
ENTRYPOINT ["/usr/sbin/vsftpd"]
EXPOSE 20 21
