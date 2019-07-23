#!/usr/bin/env bash

#cd ./../vendor/gocv.io/x/gocv/
docker build -t gocv .

#cd ./../../../../deployment
docker-compose up -d
