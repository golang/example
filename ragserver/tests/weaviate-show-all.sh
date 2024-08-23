#!/bin/bash

set -eux

echo '{
  "query": "{
    Get {
      Document { 
        text
      }
    }
  }"
}' | tr -d "\n" | curl \
    -X POST \
    -H 'Content-Type: application/json' \
    -d @- \
    http://localhost:9035/v1/graphql | jq .
