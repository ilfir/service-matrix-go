#!/bin/bash

# Define image name
IMAGE_NAME="service-matrix-go"
CONTAINER_NAME="service-matrix-go"

echo "Building Docker image..."
docker build -t $IMAGE_NAME .

# Check if container exists and remove it
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "Removing existing container..."
    docker rm -f $CONTAINER_NAME
fi

echo "Running Docker container..."
docker run -d -p 8080:8080 --name $CONTAINER_NAME $IMAGE_NAME

echo "Container started. Logs:"
docker logs -f $CONTAINER_NAME
