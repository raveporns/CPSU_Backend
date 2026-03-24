#!/bin/sh

minio server /data --console-address ":9001" &
sleep 5

mc alias set local http://localhost:9000 minioadmin minioadmin123

mc mb local/images
mc anonymous set public local/images

wait