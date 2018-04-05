FROM centos:7

# Golang Version
ENV VERSION 1.7
ENV FILE go$VERSION.linux-amd64.tar.gz

# Golang download URL
ENV URL https://dl.google.com/go/$FILE
ENV GOPATH /golang
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV GO15VENDOREXPERIMENT 1

# Adding the installation script
ADD assets/install.sh /tmp/install.sh
RUN cd /tmp/ && chmod +x install.sh && ./install.sh

# Downloading the database file from ftp
ADD assets/ftp.sh /geo-db/ftp.sh
RUN cd /geo-db/ && chmod +x ftp.sh

# This is for getting the mmdb file from the ftp server
RUN ./geo-db/ftp.sh

EXPOSE 80

CMD ["/bin/location", "serve", "--log-level=debug", "--http-addr=:80", "--db=/geo-db/city.mmdb", "--db=/geo-db/isp.mmdb", "--db=/geo-db/conn.mmdb"]
