#!/bin/bash
grpcurl -plaintext -d '{"id" : "testdevice4"}' localhost:8080 infinimesh.api.Shadows/Get | jq
