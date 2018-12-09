#!/bin/bash

echo "{\"sampleKey\" : $RANDOM}" | go run hack/mqtt_testclients/produce/produce.go -topic "shadows/testdevice3" -broker localhost:8089
