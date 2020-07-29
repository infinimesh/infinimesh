#!/bin/sh
clear

# create infinimesh directory
mkdir -p ~/infinimesh-local && cd ~/infinimesh-local

# check if we on Linux or OS X
 if [[ "$OSTYPE" == "linux-gnu" ]]; then
	 echo "Linux OS found, proceed ..."
     elif [[ "$OSTYPE" == "darwin"* ]]; then
         echo " OS X found, proceed ..."
     else
	 echo "no Linux / OSX environment, aborting ..."
       	 exit 0
fi

# check if brew ist installed
if which brew >/dev/null; then
        echo " brew is available, proceed .."
    else
        echo " brew not found, please install:" \
              " https://brew.sh/"
        exit 0
 fi

# check if kubectl is installed
if which kubectl >/dev/null; then
        echo " kubectl found, proceed .."
    else
	echo " kubectl not found, please install:" \
	echo " brew install kubectl" \
        exit 0
 fi

# check if virtualbox is installed
if which virtualbox >/dev/null; then
        echo " virtualbox found, proceed .."
    else
	echo " virtualbox not found, please install:" \
	echo " brew cask install virtualbox" \
        exit 0
 fi

# check if minikube is installed
if which minikube >/dev/null; then
        echo " minikube found, all well .."
    else
	echo " minikube not found, installing:" \
	echo " brew install minikube" \
    brew install minikube
 fi

# set virtualbox as default driver
minikube config set driver virtualbox

# start minikube
echo " we start up a small k8s with 4 CPU and 8GB RAM"
minikube stop && minikube delete
minikube start --cpus 4 --memory 8192 --disk-size 150GB

# configure minikube
minikube addons enable ingress
minikube addons enable ingress-dns

# check if we see nodes
 kubectl get nodes

# certificates
echo " creating self - signed certificates "
printf '\n'
mkdir -p ~/certs && cd ~/certs
openssl genrsa -out ca.key 4096
openssl req -subj '/CN=infinimesh.minikube/O=Infinimesh' -new -x509 -sha256 -key ca.key -out ca.crt -days 3650
openssl genrsa -out apiserver_grpc.key 4096
openssl req -subj /CN=grpc.api.infinimesh.minikube -out apiserver_grpc.csr -key apiserver_grpc.key -new
openssl x509 -req -days 3650 -in apiserver_grpc.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out apiserver_grpc.crt -sha256
openssl genrsa -out apiserver_rest.key 4096
openssl req -subj /CN=api.infinimesh.minikube -out apiserver_rest.csr -key apiserver_rest.key -new
openssl x509 -req -days 3650 -in apiserver_rest.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out apiserver_rest.crt -sha256 
openssl genrsa -out mqtt_bridge.key 4096
openssl req -subj /CN=mqtt.api.infinimesh.minikube -out mqtt_bridge.csr -key mqtt_bridge.key -new
openssl x509 -req -days 3650 -in mqtt_bridge.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out mqtt_bridge.crt -sha256 
openssl genrsa -out app.key 4096
openssl req -subj /CN=app.infinimesh.minikube -out app.csr -key app.key -new
openssl x509 -req -days 3650 -in app.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out app.crt -sha256
sleep 3

echo " install the certificates"
printf '\n'
kubectl create secret tls apiserver-grpc-tls --cert apiserver_grpc.crt --key apiserver_grpc.key 
kubectl create secret tls apiserver-rest-tls --cert apiserver_rest.crt --key apiserver_rest.key 
kubectl create secret tls mqtt-bridge-tls --cert mqtt_bridge.crt --key mqtt_bridge.key 
kubectl create secret tls app-tls --cert app.crt --key app.key 
cd -

 # setup infinimesh
echo " we setup infinimesh" \
echo " installing infinimesh operator "
printf '\n'
kubectl apply -f https://raw.githubusercontent.com/infinimesh/operator/master/manifests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/infinimesh/operator/master/manifests/operator.yaml
sleep 2

kubectl create namespace kafka
curl -L https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.14.0/strimzi-cluster-operator-0.14.0.yaml \
  | sed 's/namespace: .*/namespace: kafka/' \
  | kubectl apply -f - -n kafka 

sleep 20

kubectl apply -f https://raw.githubusercontent.com/InfiniteDevices/infinimesh/master/hack/minikube/infinimesh-kafka.yaml -n kafka
sleep 20

kubectl apply -f https://raw.githubusercontent.com/InfiniteDevices/infinimesh/master/hack/minikube/infinimesh-platform.yaml
sleep 20

echo "=> installing inf (infinimesh CLI) and point to the local setup:"
curl -L https://bit.ly/2CNKWzJ | BINDIR=$HOME/bin bash  
echo "inf CLI installed in" $HOME/bin
~/bin/inf config set-context local --apiserver grpc.api.infinimesh.minikube:443 --tls=true --ca-file ~/certs/ca.crt

