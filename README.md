# traefik-operator
This repository contains the source code of an operator for orchestrating Traefik.
The project is resulted by a design research due to a master thesis.

## Design
### 1. Project Setup
```shell
kubebuilder init --domain mh.edu.com --repo github.com/Education-Orga/traefik-operator
```

### 2. API Creation
```shell
# kind TraefikInstance
kubebuilder create api --group traefik --version v1alpha1 --kind TraefikInstance
```

### 3. Business Logic
- CR Development of kind TraefikInstance
- Custom Controller Logic /pkg/deployment/deploy.go
- Reconcile: integration of business logic /internal/traefikinstance_controller.go

## Testing
xxx

## Run the operator
### 1. Local outside of a cluster
Clone the repository and run:
```shell
make run
```
### 2. Docker deployment inside of a cluster
Clone the repository and run:
```shell
# build the docker image 
make docker-build IMG=<image-name>

# push the docker image 
make docker-push IMG=<image-name>

# deploy the operator to the cluster
make deploy

# undeploy the operator to the cluster
make undeploy
```