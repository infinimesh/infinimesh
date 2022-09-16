#!/bin/sh
for image in Dockerfiles/*; do
    tag="ghcr.io/infinimesh/infinimesh/$(basename $image):latest"
    docker build . -f "$image/Dockerfile" -t $tag
    docker push $tag
done