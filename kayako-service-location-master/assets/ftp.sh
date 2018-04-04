ftp -inpv gbuild.gdev.com << EOF
user gbuildftp EALhMiuJ8Ry
cd /ftp/gbuildftp/kayako/location/
get city.mmdb
get isp.mmdb
get conn.mmdb
bye
EOF
