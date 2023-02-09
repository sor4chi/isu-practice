DIR_PATH=$(cd $(dirname $0); pwd)
cat $DIR_PATH/../logs/nginx/access.log | alp json
