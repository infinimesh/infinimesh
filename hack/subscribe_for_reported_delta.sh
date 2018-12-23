#!/bin/bash 

grpcurl -plaintext -d '{"id" : "testdevice4"}' localhost:8096 infinimesh.shadow.Shadows/StreamReportedStateChanges
