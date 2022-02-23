sec=5
name=xamppBCC
mysqlpassword=123

echo "Creating container with the name $name"; podman run --name $name -p 41061:22 -p 41062:80 -p 41063:3306 -p 41064:5900 -d -v /www:/www tomsik68/xampp
echo ""
while [ $sec -ge 0 ]; do echo -ne "Please wait for $sec seconds\033[0K\r"; let "sec=sec-1"; sleep 1; done
echo ""
echo ""
echo "Adding PATH to .bashrc"; podman exec -it $name perl -0644 -i -pe "s/# alias mv\=\'mv \-i\'/# alias mv\=\'mv \-i\'\nexport PATH\=\/opt\/lampp\/bin\:\\\$PATH/igs" /root/.bashrc
echo ""
echo "creating default database for bcc_backend"; podman exec -it $name /opt/lampp/bin/mysql -u root -e 'CREATE DATABASE bcc_backend;'
echo "Adding access to MySQL"; podman exec -it $name /opt/lampp/bin/mysql -u root -D mysql -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '$mysqlpassword' WITH GRANT OPTION;"
echo ""
echo "Moving document root into /www"; podman exec -it $name perl -0644 -i -pe 's/DocumentRoot \"\/opt\/lampp\/htdocs\"/#DocumentRoot \"\/opt\/lampp\/htdocs\"\nDocumentRoot \"\/www\"/;s/\<Directory \"\/opt\/lampp\/htdocs\"\>/#\<Directory \"\/opt\/lampp\/htdocs\"\>\n\<Directory \"\/www\">/igs' /opt/lampp/etc/httpd.conf
echo ""
echo "Restarting $name"; podman restart $name
echo ""
echo "Please connect through these ports"; podman port $name
