#!/bin/sh

cp /gophermarttest ./gophermarttest
echo "Build gophermart"
cd cmd/gophermart && go build -buildvcs=false -o gophermart && cd ../../

echo "Start command"
echo "$@"
exec "$@"