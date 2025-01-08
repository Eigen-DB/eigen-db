#!/bin/bash

which helm
if [ $? != 0 ]; then
    echo "Helm is not installed on your system.";
    exit 1
fi

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace
