#!/bin/sh
version="${1:-latest}"
for image in Dockerfiles/*; do
    tag="ghcr.io/infinimesh/infinimesh/$(basename $image):latest"
    docker build . -f "$image/Dockerfile" -t $tag
done