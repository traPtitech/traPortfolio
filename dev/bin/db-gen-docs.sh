#!/bin/bash

set -eu

if [ -d "./docs/dbschema" ]; then
    rm -r ./docs/dbschema
fi
DB_HOST=localhost go run main.go &

docker run --rm --net=host -e TBLS_DSN="mysql://root:password@127.0.0.1:${1}/portfolio" -v $(pwd):/work k1low/tbls:${2} doc
kill $(ps | grep "go run main.go" | grep -v grep | awk '{print $1}')
