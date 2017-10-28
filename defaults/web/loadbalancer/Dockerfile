FROM alpine
RUN apk update
RUN apk add nginx
RUN mkdir -p /run/nginx
RUN sed '2 i daemon off;' -i /etc/nginx/nginx.conf
ENTRYPOINT ["/usr/sbin/nginx"]
EXPOSE 80 443