FROM golang:1.5-alpine
RUN apk update
RUN apk add git
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/America/Denver /etc/localtime
ADD root /var/spool/cron/crontabs/root
RUN mkdir -p /go/src/oct-memcached-preprovision
ADD oct-memcached-preprovision.go  /go/src/oct-memcached-preprovision/oct-memcached-preprovision.go
ADD build.sh /build.sh
RUN chmod +x /build.sh
RUN /build.sh
#CMD ["/go/src/oct-memcached-preprovision/oct-memcached-preprovision"]
CMD ["crond", "-f"]
#RUN mkdir /root/.aws
#ADD credentials /root/.aws/credentials




