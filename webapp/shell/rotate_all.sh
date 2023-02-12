FILE_DIR=$(
    cd $(dirname $0)
    pwd
)

sh $FILE_DIR/nginx_rotation.sh
sh $FILE_DIR/mysql_rotation.sh
