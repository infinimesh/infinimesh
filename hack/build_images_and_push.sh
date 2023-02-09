#!/bin/sh
for image in Dockerfiles/*; do
    tag="ghcr.io/infinimesh/infinimesh/$(basename $image):latest"
    INFINIMESH_VERSION_TAG=$(git describe --tags --abbrev=0) docker build . -f "$image/Dockerfile" -t $tag
    docker push $tag
done