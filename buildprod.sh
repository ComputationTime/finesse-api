#!/bin/bash
docker compose -f docker-compose-prod.yml build
docker push computationtime/finesse-api:amd64
ssh finesse << EOF
# sudo docker stop $(sudo docker ps -q --filter ancestor=computationtime/finesse-api:amd64)
sudo docker stop $(sudo docker ps -a -q)
sudo docker image rm -f computationtime/finesse-api:amd64
sudo docker pull computationtime/finesse-api:amd64
sudo docker run -v $(pwd)/finesse-frontend/build:/usr/src/api/build -p 8000:8000 -d computationtime/finesse-api:amd64;
EOF