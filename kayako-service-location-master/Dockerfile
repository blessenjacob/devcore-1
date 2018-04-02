FROM centos:7

# Golang Version
ENV VERSION 1.7
ENV FILE go$VERSION.linux-amd64.tar.gz
# Golang download URL
ENV URL https://dl.google.com/go/$FILE
ENV GOPATH /golang
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

#Installation of required packagaes
RUN yum -y install git make

# Install of golang
RUN yum -y clean all
RUN curl -OL $URL
RUN tar -C /usr/local -xzf $FILE
RUN rm -rf $FILE


RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" /geo-db && chmod -R 777 "$GOPATH"

ADD ./db-files/city.mmdb /geo-db/city.mmdb
ADD ./db-files/isp.mmdb  /geo-db/isp.mmdb
ADD ./db-files/conn.mmdb /geo-db/conn.mmdb

ENV GO15VENDOREXPERIMENT 1

ADD . $GOPATH/src/github.com/kayako/service-location
RUN cd $GOPATH/src/github.com/kayako/service-location && make geo && mv ./build/location /bin/location

RUN rm -rf $GOPATH

EXPOSE 80

ENTRYPOINT ["/bin/location", "serve", "--log-level=debug", "--http-addr=:80", "--db=/geo-db/city.mmdb", "--db=/geo-db/isp.mmdb", "--db=/geo-db/conn.mmdb"]
