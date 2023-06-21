#!/bin/bash


if [[ -n $(docker ps -q -f "name=mongodb") ]];then
	docker stop mongodb
else
	echo "container had stopped"
    exit 0
fi