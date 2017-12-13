FROM ubuntu
RUN apt-get update
RUN apt-get install --asume-yes true vsftpd
RUN echo "anon_root=/var/lib/ftp/" >> /etc/vsftpd.conf
RUN sed -i "/write_enable/s/#write/write/" /etc/vsftpd.conf
RUN echo "seccomp_sandbox=NO" >> /etc/vsftpd.conf
RUN sed -i "/xferlog_file/s/#xferlog_file/xferlog_file/g" /etc/vsftpd.conf
RUN mkdir -p /var/run/vsftpd/empty
RUN groupadd usuariosftp -g 1001
RUN echo "pam_service_name=ftp" >> /etc/vsftpd.conf
RUN echo "pasv_enable=YES" >> /etc/vsftpd.conf
RUN echo "pasv_min_port=21100" >> /etc/vsftpd.conf
RUN echo "pasv_max_port=21110" >> /etc/vsftpd.conf
RUN echo "pasv_address=35.153.46.76" >> /etc/vsftpd.conf
RUN echo "port_enable=YES" >> /etc/vsftpd.conf
RUN echo "connect_from_port_20=YES" >> /etc/vsftpd.conf
RUN echo "ftp_data_port=20" >> /etc/vsftpd.conf
ENTRYPOINT ["/usr/sbin/vsftpd"]
EXPOSE 20 21
