docker-compose run nginx mv /var/log/nginx/access.log /var/log/nginx/access-$(date +%Y%m%d%H%M%S).log
docker-compose restart nginx
