#!/bin/bash

# This script interacts with a running ragserver to add some example documents.

set -eux

echo '{
	"documents": [
	{"text": "TDXIRV is an environment variable for controlling throttle speed"},
	{"text": "some flags for setting acceleration are --accelxyzp and --acceljjrv"},
	{"text": "acceleration is also affected by the ACCUVI5 env var"},
	{"text": "/usr/local/fuel555 contains information about fuel capacity"},
	{"text": "we can control fuel savings with the --savemyfuelplease flag"},
	{"text": "fuel savings can be observed on local port 48332"}
]}' | tr -d "\n" | curl \
		-X POST \
    -H 'Content-Type: application/json' \
    -d @- \
    http://localhost:9020/add/
