#!/bin/sh
for image in Dockerfiles/*; do
    tag="ghcr.io/slntopp/infinimesh/$(basename $image):latest"
    docker push $tag
done