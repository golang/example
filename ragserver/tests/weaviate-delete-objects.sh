#!/bin/bash

set -eux

curl \
  -X DELETE \
  http://localhost:9035/v1/schema/Document
