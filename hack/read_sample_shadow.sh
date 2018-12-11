#!/bin/bash
grpcurl -plaintext -d '{"device_id" : "testdevice4"}' localhost:8080 infinimesh.api.Shadow/GetReported | jq
