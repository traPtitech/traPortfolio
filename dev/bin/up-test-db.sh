#!/bin/bash

set -eu

containername=portfolio-test-db
port=${TEST_DB_PORT:-3306}

if docker ps | grep ${containername} >/dev/null; then
    exit 0 # 既にテストDBコンテナが起動している
fi

if docker ps --all | grep ${containername} >/dev/null; then
    echo "restart ${containername} docker container"
    docker restart ${containername}
else
    echo "create ${containername} docker container"
    docker run --name ${containername} -p ${port}:3306 -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=portfolio -d mariadb:10.0.19 \
        mysqld --character-set-server=utf8 --collation-server=utf8_general_ci
fi
