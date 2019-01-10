#!/bin/bash

echo "{\"sampleKey\" : $RANDOM}" | go run hack/mqtt_testclients/produce/produce.go -topic "shadows/testdevice4" -broker localhost:8089
#echo "{\"sampleKey\" : $RANDOM}" | go run hack/mqtt_testclients/produce/produce.go -topic "shadows/testdevice4" -broker controlplane.infinimesh.io:8089
