#!/bin/bash

sudo docker stop tech-db
sudo docker rm tech-db
sudo docker build -t tech-db -f Dockerfile .
sudo docker run -p 5000:5000 --name tech-db -t tech-db
