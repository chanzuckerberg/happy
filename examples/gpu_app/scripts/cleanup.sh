#!/bin/sh

# This script is meant to be a measure of last resort if you deleted the integration-test stack with `happy delete --force` and it failed to delete all resources.
# Please do not use it otherwise. AWS profile is ommitted because it is assumed that you are running this script from the github action you ran `happy delete --force` from.

kubectl delete service,ing,secret,configmap,deployment,serviceaccount,hpa -l app.kubernetes.io/name=integration-test -n si-rdev-happy-eks-rdev-happy-env || true
kubectl delete secret integration-test-oidc-config -n si-rdev-happy-eks-rdev-happy-env || true
kubectl delete svc integration-test-frontend -n si-rdev-happy-eks-rdev-happy-env || true
kubectl delete ing integration-test-frontend -n si-rdev-happy-eks-rdev-happy-env || true
kubectl delete ing integration-test-frontend-options-bypass -n si-rdev-happy-eks-rdev-happy-env || true
kubectl delete serviceaccount integration-test-frontend-rdev-integration-test -n si-rdev-happy-eks-rdev-happy-env || true
kubectl delete deployment integration-test-frontend -n si-rdev-happy-eks-rdev-happy-env || true
aws iam delete-role --role-name integration-test-frontend-rdev-integration-test --region us-west-2 || true