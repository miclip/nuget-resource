#!/bin/bash

set -ex 

docker build . -t nuget-resource -t miclip/nuget-resource:latest
docker push miclip/nuget-resource:latest