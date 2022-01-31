````
docker build -t pow-app -f docker/Dockerfile.App .

docker build -t pow-client -f docker/Dockerfile.Client .

docker network create pow-net --driver bridge

docker run -p 1337:1337 --network pow-net \
--name pow-app -itd \
-e COUNT_OF_ZERO_BYTES=1 \
-e HTTP_PORT=1337 \
-e MAC_SECRET_KEY=somekey \
-e CHALLENGE_TTL=1m \
pow-app

docker run --network pow-net --name pow-client \ 
-e QUOTE_HOST=http://pow-app:1337 \
pow-client

docker start -i pow-client
````
