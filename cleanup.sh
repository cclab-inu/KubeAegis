#!/bin/bash

# Delete all KubeAegis resources
kubectl delete KubeAegisPolicy --all --all-namespaces

# Delete all KubeArmorPolicy resouces
kubectl delete cnp --all --all-namespaces
kubectl delete ksp --all --all-namespaces
kubectl delete kyverno --all --all-namespaces

echo "All resources have been successfully deleted."