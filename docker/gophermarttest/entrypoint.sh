#!/bin/sh

cp /gophermarttest ./gophermarttest
echo "Build gophermart"
cd cmd/gophermart && go build -buildvcs=false -o gophermart && cd ../../

echo $DATABASE_URI


echo "Start command"
echo "$@"
exec "$@"