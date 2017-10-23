FROM alpine:latest
RUN apk update
RUN apk add bind
RUN mkdir /etc/bind/zones/
RUN mkdir /etc/bind/conf.d/
CMD /usr/sbin/named -g
EXPOSE 53 53/udp
