#!/bin/bash
set -ex

cd ~/Youyu-se/
git fetch && git rebase origin/master 
docker-compose -f docker/docker-compose.yml restart 

