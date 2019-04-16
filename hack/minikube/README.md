# Setup instructions for minikube

## Installing operators
Infinimesh is a Cloud Native Application and leverages Kubernetes Operators to install and maintain the platform. Depending on your setup, you may need to install:
- Infinimesh Operator
- Kafka Strimzi Operator (optional)
- KubeDB Operator

### Infinimesh Operator
The following commands install the Infinimesh Kubernetes Operator into `infinimesh-system`.
```
kubectl apply -f https://raw.githubusercontent.com/infinimesh/operator/master/manifests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/infinimesh/operator/master/manifests/operator.yaml
```

### KubeDB Operator
To provision Postgres and Redis instances, infinimesh uses `KubeDB`.
```
curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0/hack/deploy/kubedb.sh | bash
```

### Strimzi Kafka Operator (optional)
Strimzi is a Kubernetes operator to install and maintain Kafka in a Cloud Native way. Strimzi is not required by Infinimesh. However, if you don't already have a Kafka installation, it is the recommended way to install Kafka.
*If you already have a Kafka Cluster, you can skip this step.*

Install Strimzi Kafka Operator into the namespace `kafka`:
```
kubectl create namespace kafka
curl -L https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.11.1/strimzi-cluster-operator-0.11.1.yaml \
  | sed 's/namespace: .*/namespace: kafka/' \
  | kubectl -n kafka apply -f -
```

For more details refer to [![Strimzi Documentation](https://strimzi.io/quickstarts/minikube/)].

### NGINX Ingress Controller

For API Servers and App, infinimesh uses [![Nginx Ingress Controller](https://kubernetes.github.io/ingress-nginx/)]. When using Minikube, enable it as addon:
```
minikube addons enable ingress
```

When deploy to a real cluster, follow the [![instructions](https://kubernetes.github.io/ingress-nginx/deploy/)]. This involves installing some manifests or a HELM chart and may vary slightly, depending on your infrastructure.

## TLS Certicate Setup
Infinimesh requires X509 KeyPairs for the following applications:
- API Server gRPC
- API Server REST
- MQTT Bridge
- App (User Interface)

To manage the certificates, infinimesh uses `cert-manager`. It is a Kubernetes operator to provision certificates. A typical usage is provisioning via `Let's Encrypt`, but infinimesh uses it for all kinds of certificates.
TODO

You may serve these domains from any host/domain you want. It's just critical that you have TLS server certificates for each of them, with the used domain name.
For local installation, certificates can be generated with `openssl` and self-signed. This is not recommended for producation scenarios, because other parties (e.g. users of the platform) will have to import the certificates into their trust store. This is not just inconvenient, but in most scenarios also makes the platform vulnerable to `Man in the middle atacks`. Thus, use self-signed certificates only local & testing environments.

You need `openssl` for the following steps.

0. Generate Root CA
0.1 Generate Private Key
```
openssl genrsa -out ca.key 4096
openssl req -subj / -new -x509 -sha256 -key ca.key -out ca.crt -days 3650
```

1. Generate Platform Private Keys & Certificates
1.1 API Server
```
openssl genrsa -out apiserver_grpc.key 4096
openssl req -subj /CN=grpc.api.infinimesh.localhost -out apiserver_grpc.csr -key apiserver_grpc.key -new
openssl x509 -req -days 3650 -in apiserver_grpc.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out apiserver_grpc.crt -sha256 
```
1.2 API Server REST
```
openssl genrsa -out apiserver_rest.key 4096
openssl req -subj /CN=api.infinimesh.localhost -out apiserver_rest.csr -key apiserver_rest.key -new
openssl x509 -req -days 3650 -in apiserver_rest.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out apiserver_rest.crt -sha256 
```
1.3 MQTT Bridge
```
openssl genrsa -out mqtt_bridge.key 4096
openssl req -subj /CN=mqtt.api.infinimesh.localhost -out mqtt_bridge.csr -key mqtt_bridge.key -new
openssl x509 -req -days 3650 -in mqtt_bridge.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out mqtt_bridge.crt -sha256 
```

1.4 App
```
openssl genrsa -out app.key 4096
openssl req -subj /CN=app.infinimesh.localhost -out app.csr -key app.key -new
openssl x509 -req -days 3650 -in app.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out app.crt -sha256 
```

If you want to use different domains, replace `localhost` in the `-subj` parameter.

Now we transform these Cert & Key files into Kubernetes secrets, and deploy them to the cluster.

```
kubectl create secret tls apiserver-grpc-tls --cert apiserver_grpc.crt --key apiserver_grpc.key 
kubectl create secret tls apiserver-rest-tls --cert apiserver_rest.crt --key apiserver_rest.key 
kubectl create secret tls mqtt-bridge-tls --cert mqtt_bridge.crt --key mqtt_bridge.key 
kubectl create secret tls app-tls --cert app.crt --key app.key 
```

## Deploy Kubernetes Resources
Now we're installing the resources required by Infinimesh into Kubernetes:
- Infinimesh itself via the `Platform` CRD
- Kafka via Strimzi Operator (optional)
- Databases

Note: We already installed Secrets in the previous step.

1. Infinimesh Platform Resource

```
apiVersion: infinimesh.infinimesh.io/v1beta1
kind: Platform
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: my-platform
spec:
  kafka:
    bootstrapServers: "infinimesh-kafka-bootstrap.kafka.svc.cluster.local:9092"
  app:
    host: "app.infinimesh.localhost"
    tls:
      - hosts:
        - "app.infinimesh.localhost"
        secretName: "app-tls"
  mqtt:
    secretName: "mqtt-bridge-tls"
  apiserver:
    restful:
      host: "api.infinimesh.localhost"
      tls:
        - hosts:
          - "api.infinimesh.localhost"
          secretName: "apiserver-rest-tls"
    grpc:
      host: "grpc.api.infinimesh.localhost"
      tls:
        - hosts:
          - "grpc.api.infinimesh.localhost"
          secretName: "apiserver-grpc-tls"
```

Take care that the host&tls config match to the deployed secrets.

2. Kafka
```
apiVersion: kafka.strimzi.io/v1alpha1
kind: Kafka
metadata:
  name: infinimesh
spec:
  kafka:
    version: 2.1.0
    replicas: 3
    listeners:
      plain: {}
      tls: {}
    config:
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      log.message.format.version: "2.1"
    storage:
      type: persistent-claim
      size: 100Gi
      deleteClaim: false
  zookeeper:
    replicas: 3
    storage:
      type: persistent-claim
      size: 100Gi
      deleteClaim: false
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

3. NGINX Ingress Controller
TODO

4. TimescaleDB
TODO
