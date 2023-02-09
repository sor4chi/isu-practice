docker-compose run mysql mv /var/log/mysql/mysql-slow.log /var/log/mysql/mysql-slow-$(date +%Y%m%d%H%M%S).log
docker-compose restart mysql
