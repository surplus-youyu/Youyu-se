#!/bin/bash
set -ex

cd ~/Youyu-se/
git pull
docker-compose -f docker/docker-compose.yml up -d

