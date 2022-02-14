#!/bin/sh
for image in Dockerfiles/*; do
    tag="ghcr.io/slntopp/infinimesh/$(basename $image):latest"
    docker build . -f "$image/Dockerfile" -t $tag
    docker push $tag
done