
    [vagrant@localhost docker-q8sad]$ docker build -t tangfeixiong/pd-server
    docker: "build" requires 1 argument.
    See 'docker build --help'.
    
    Usage:	docker build [OPTIONS] PATH | URL | -
    
    Build an image from a Dockerfile
    [vagrant@localhost docker-q8sad]$ docker build -t tangfeixiong/pd-server .
    Sending build context to Docker daemon 19.73 MB
    Step 1 : FROM alpine:3.4
     ---> 31b45a1205be
    Step 2 : MAINTAINER tangfeixiong <fxtang@qingyuanos.com>
     ---> Using cache
     ---> 5725fae9db2b
    Step 3 : RUN apk add --update bash ca-certificates
     ---> Running in 5f342c819cfb
    fetch http://dl-cdn.alpinelinux.org/alpine/v3.4/main/x86_64/APKINDEX.tar.gz
    fetch http://dl-cdn.alpinelinux.org/alpine/v3.4/community/x86_64/APKINDEX.tar.gz
    (1/6) Installing ncurses-terminfo-base (6.0-r7)
    (2/6) Installing ncurses-terminfo (6.0-r7)
    (3/6) Installing ncurses-libs (6.0-r7)
    (4/6) Installing readline (6.3.008-r4)
    (5/6) Installing bash (4.3.42-r3)
    Executing bash-4.3.42-r3.post-install
    (6/6) Installing ca-certificates (20160104-r4)
    Executing busybox-1.24.2-r9.trigger
    Executing ca-certificates-20160104-r4.trigger
    OK: 14 MiB in 17 packages
     ---> 0a81f8ef7fd7
    Removing intermediate container 5f342c819cfb
    Step 4 : ADD pd-server /pd-server
     ---> c2080385d2ec
    Removing intermediate container 1d18a4b3b5ca
    Step 5 : ENV Q8S_ETCD_ENV ""
     ---> Running in 06e5d6ce910e
     ---> b7afcec40b99
    Removing intermediate container 06e5d6ce910e
    Step 6 : EXPOSE 1234
     ---> Running in 4aefd87114e2
     ---> 1ed9a7314fe4
    Removing intermediate container 4aefd87114e2
    Step 7 : ENTRYPOINT /pd-server
     ---> Running in 27ca4346b8af
     ---> fba04dc5cb03
    Removing intermediate container 27ca4346b8af
    Successfully built fba04dc5cb03
