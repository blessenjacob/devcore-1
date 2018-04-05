ftp -inpv gbuild.gdev.com << EOF
user gbuildftp EALhMiuJ8Ry
cd /ftp/gbuildftp/kayako/location/
bin
get city.mmdb /geo-db/city.mmdb
get isp.mmdb  /geo-db/isp.mmdb
get conn.mmdb /geo-db/conn.mmdb
get vendor.tar.gz /golang/src/github.com/kayako/service-location/vendor.tar.gz
bye
EOF

cd /golang/src/github.com/kayako/service-location/ && tar -xvf vendor.tar.gz  &&  make geo && mv ./build/location /bin/location
rm -rf /golang
