#!/bin/bash
# Note: You must first create the ssl-sidecar registry in the AWS account ECR (for czi-playground it already exists)
tag="0.0.6"
profile="czi-playground"
account=${ACCOUNTID}
aws ecr get-login-password --profile $profile | docker login --username AWS --password-stdin $account.dkr.ecr.us-west-2.amazonaws.com
docker build . --no-cache --platform linux/arm64 -t $account.dkr.ecr.us-west-2.amazonaws.com/ssl-sidecar:$tag
docker push $account.dkr.ecr.us-west-2.amazonaws.com/ssl-sidecar:$tag