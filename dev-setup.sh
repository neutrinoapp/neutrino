docker -v

docker pull apcera/gnatsd
docker pull mongo:latest

docker run --name nats -p 4222:4222 -p 8333:8333 -d apcera/gnatsd -m 8333
docker run --name mongodb -p 27017:27017 -v /data/mongo:/data/db -d mongo --storageEngine wiredTiger

docker ps

echo "All done!"
