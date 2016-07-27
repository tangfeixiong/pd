FROM alpine:3.4
MAINTAINER tangfeixiong <fxtang@qingyuanos.com>

RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*

ADD pd-server /pd-server

ENV Q8S_ETCD_ENV=""

EXPOSE 1234

CMD ["/pd-server"]

# ADD run.sh /run.sh
# CMD ["/run.sh"]
