#!/bin/bash

# Set the URL to make the HTTP request to
URL="http://example.com/api/endpoint"

# Set the JSON payload to send in the request
JSON_PAYLOAD='{"key1": "value1", "key2": "value2"}'

# Make the HTTP request using curl
curl --request POST \
     --header "Content-Type: application/json" \
     --data "$JSON_PAYLOAD" \
     "$URL"
