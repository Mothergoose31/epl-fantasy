=#!/bin/bash

set -e  

NETWORK_NAME="epl-fantasy-network"
MONGODB_CONTAINER_NAME="mongodb"
APP_CONTAINER_NAME="epl-fantasy"

echo "Checking network..."
if ! docker network inspect $NETWORK_NAME >/dev/null 2>&1; then
    echo "Creating network '$NETWORK_NAME'..."
    docker network create $NETWORK_NAME
fi

echo "Checking MongoDB container..."
if ! docker ps -a --format '{{.Names}}' | grep -q "^$MONGODB_CONTAINER_NAME$"; then
    echo "Creating MongoDB container..."
    docker run --name $MONGODB_CONTAINER_NAME --network $NETWORK_NAME -p 27017:27017 -d mongo
else
    echo "MongoDB container exists. Ensuring it's running and connected to the network..."
    docker start $MONGODB_CONTAINER_NAME
    if ! docker network inspect $NETWORK_NAME | grep -q "\"$MONGODB_CONTAINER_NAME\""; then
        docker network connect $NETWORK_NAME $MONGODB_CONTAINER_NAME
    fi
fi

echo "Waiting for MongoDB to be ready..."
until docker exec $MONGODB_CONTAINER_NAME mongosh --eval "db.runCommand('ping').ok" --quiet >/dev/null 2>&1
do
    echo "Waiting for MongoDB to be ready..."
    sleep 2
done
echo "MongoDB is up and running"

echo "Building epl-fantasy Docker image..."
docker build -t $APP_CONTAINER_NAME .

echo "Checking epl-fantasy container..."
if docker ps -a --format '{{.Names}}' | grep -q "^$APP_CONTAINER_NAME$"; then
    echo "Removing existing epl-fantasy container..."
    docker rm -f $APP_CONTAINER_NAME
fi

echo "Creating and starting new epl-fantasy container..."
docker run --name $APP_CONTAINER_NAME --network $NETWORK_NAME -p 8080:8080 -d $APP_CONTAINER_NAME

echo "Verifying all containers are in the network..."
for container in $MONGODB_CONTAINER_NAME $APP_CONTAINER_NAME; do
    if ! docker network inspect $NETWORK_NAME | grep -q "\"$container\""; then
        echo "Connecting $container to $NETWORK_NAME..."
        docker network connect $NETWORK_NAME $container
    fi
done

echo "Network has been created, all containers are up and running"