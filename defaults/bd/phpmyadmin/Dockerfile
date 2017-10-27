FROM netroby/alpine-php

ENV PMA_VER 4_6_5_2
RUN wget --no-check-certificate -c https://github.com/phpmyadmin/phpmyadmin/archive/RELEASE_$PMA_VER.tar.gz; \
        tar zxvf RELEASE_$PMA_VER.tar.gz;\
        mv phpmyadmin-RELEASE_$PMA_VER /www; \
        unlink *.tar.gz
RUN sed -i -e "s/^upload_max_filesize\s*=\s*2M/upload_max_filesize = 64M/" /etc/php5/php.ini
RUN sed -i -e "s/^post_max_size\s*=\s*8M/post_max_size = 64M/" /etc/php5/php.ini

COPY config.inc.php /www/config.inc.php

EXPOSE 8080

ENTRYPOINT ["php", "-S", "0.0.0.0:8080", "-t", "/www/"]

