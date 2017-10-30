FROM alpine
RUN apk update
RUN apk add vsftpd
RUN echo "anon_root=/var/lib/ftp/" >> /etc/vsftpd/vsftpd.conf
RUN sed -i "/write_enable/s/#write/write/" /etc/vsftpd/vsftpd.conf
RUN echo "seccomp_sandbox=NO" >> /etc/vsftpd/vsftpd.conf
RUN sed -i "/xferlog_file/s/#xferlog_file/xferlog_file/g" /etc/vsftpd/vsftpd.conf
ENTRYPOINT ["/usr/sbin/vsftpd"]
EXPOSE 20 21
