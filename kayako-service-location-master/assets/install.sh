#Installation of required packagaes
yum -y install git make ftp
yum -y clean all
curl -OL $URL
tar -C /usr/local -xzf $FILE
rm -rf $FILE

mkdir -p "$GOPATH/src/github.com/kayako/service-location" "$GOPATH/bin" /geo-db && chmod -R 777 "$GOPATH"
