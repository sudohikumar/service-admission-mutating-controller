#!/usr/bin/env bash

# Setup
export SERVICE_NAME=service-webhook
export NAMESPACE=admission-controller
export KEY_NAME=server.key
export CRT_NAME=server.crt
export CA_KEY=ca.key
export CA_CRT=ca.crt

# Create Namespace
kubectl create namespace $NAMESPACE

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout $CA_KEY -out $CA_CRT -subj "/CN=Service Admission Controller CA"

# Generate the private key for the webhook server
openssl genrsa -out $KEY_NAME 2048

# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key $KEY_NAME -subj "/CN=$SERVICE_NAME.$NAMESPACE.svc" \
    | openssl x509 -req -CA $CA_CRT -CAkey $CA_KEY -CAcreateserial -out $CRT_NAME

ca_pem_b64="$(openssl base64 -A <$CA_CRT)"
echo "$ca_pem_b64"

# Create secret
kubectl -n $NAMESPACE create secret tls admission-controller-tls \
    --cert $CRT_NAME \
    --key $KEY_NAME