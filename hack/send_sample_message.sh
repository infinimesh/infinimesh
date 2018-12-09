#!/bin/bash

echo '{"sampleKey" : "sampleString"}' | go run hack/mqtt_testclients/produce/produce.go -topic "shadows/testdevice3" -broker localhost:8089
