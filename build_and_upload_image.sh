#!/bin/bash


DOCKERFILE_PATH=$1
IMAGE_NAME=$2

if [ "$DOCKERFILE_PATH" = "" ]; then
    echo "Positional parameter 1 - DOCKERFILE_PATH is empty"
    exit 1
fi

if [ "$IMAGE_NAME" = "" ]; then
    echo "Positional parameter 2 - IMAGE_NAME is empty"
    exit 1
fi

echo "DOCKERFILE_PATH=" $DOCKERFILE_PATH
echo "IMAGE_NAME=" $IMAGE_NAME


docker build -f $DOCKERFILE_PATH -t $IMAGE_NAME .
echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin

IMAGE_TAG=$(date "+%Y-%m-%dT%H-%M-%S")
docker tag $IMAGE_NAME:latest $IMAGE_NAME:$IMAGE_TAG
docker push $IMAGE_NAME:latest
docker push $IMAGE_NAME:$IMAGE_TAG
