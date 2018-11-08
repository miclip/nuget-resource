#!/bin/bash

set -ex 

docker build . -t nuget-resource -t localhost:5000/nuget-resource:latest
docker push localhost:5000/nuget-resource:latest