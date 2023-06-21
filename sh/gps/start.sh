#!/bin/bash

# shell 中利用 -n 来判定字符串非空
if [[ -n $(docker ps -q -f "name=mongodb") ]];then
	echo "container is running"
    exit 0
else
	docker run -d --rm --privileged --net=host --name mongodb mongo:4.0 --auth
fi