ftp -inpv gbuild.gdev.com << EOF
user gbuildftp EALhMiuJ8Ry
cd /ftp/gbuildftp/kayako/location/
get city.mmdb /geo-db/city.mmdb
get isp.mmdb  /geo-db/isp.mmdb
get conn.mmdb /geo-db/conn.mmdb
bye
EOF
