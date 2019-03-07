#!/bin/bash

export TOKEN=$(curl -s -X POST -d '{"username" : "joe", "password": "test123"}'  localhost:8081/accounts/token | jq -r ".token")
echo "{\"sampleKey\" : $RANDOM}" | go run hack/mqtt_testclients/produce/produce.go -topic "shadows/testdevice4" -broker localhost:8089
