#!/bin/bash

export TOKEN=$(curl -s -X POST -d '{"username" : "joe", "password": "test123"}'  localhost:8081/token | jq -r ".token")
grpcurl -H "authorization: bearer $TOKEN" -plaintext -d "{\"id\" : \"testdevice4\", \"data\" : $RANDOM }" localhost:8000 infinimesh.api.Shadows/PatchDesiredState
