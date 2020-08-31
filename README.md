## Admission Controller & Mutating Webhook

This is a simple admission controller written in golang. Objective of the controller is that this controller intercepts incoming `k8s services` and checks if the name of the service has `simple` in it. 
If the service contains `simple`, this controller rejects that service and gives error as follows.

```bash
Error from server: admission webhook "service-webhook.admission-controller.svc" denied the request: service name should not contain 'simple' word
```

If the service doesn't have `simple` word, it will mutate the service to add the following label
 
```
mutated-via-controller: "true"
```

### Setup

**Deployment Name**: service-admission-controller

**Service Name**: service-webhook

**Controller Name**: service-webhook

**Namespace**: admission-controller

Run `setup.sh` file to generate certs and add it to secrets in k8s.

### Generate Certs

```bash
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
```

### How to run

- Create the certificate and put as secret in kubernetes (Refer `setup.sh` file). All the files are stored in `secrets` directory.
- Create go binary as follows
 
```bash
CGO_ENABLED=0 GOOS=linux go build -o image/service-controller
```

- This will put `service-controller` binary in `image` folder.
- Navigate to `image` folder and build docker image (here image name is `service-controller`)

```bash
docker build -t service-controller .
```

- Navigate to deployment folder and do 

```bash
kubectl apply -f deployment.yaml
```

### NOTE

- Replace `$CA_BUNDLE_VALUE` in `deployment.yaml` to the value obtained from `setup.sh` script. Value is stored in variable `$ca_pem_b64` in script.
- Secret is mounted in `deployment.yaml` on the path `/secrets/tls`. This path is used in code to do TLS Authentication. 
- Admission Controllers only run with TLS configuration.

### References:
- https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers
