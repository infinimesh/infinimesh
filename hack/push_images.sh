#!/bin/sh
for image in Dockerfiles/*; do
    tag="ghcr.io/infinimesh/infinimesh/$(basename $image):latest"
    docker push $tag
done