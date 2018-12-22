#!/bin/bash

grpcurl -plaintext -d "{\"id\" : \"testdevice4\", \"data\" : $RANDOM }" localhost:8080 infinimesh.api.Shadows/PatchDesiredState
