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

# check if multipass is installed
 if which multipass >/dev/null; then
        echo " multipass found, proceed .."
    else
        echo " multipass not found, please install:" \
              " https://github.com/CanonicalLtd/multipass/releases"
        exit 0
 fi

# check if kubectl is installed
if which kubectl >/dev/null; then
        echo " kubectl found, proceed .."
    else
	echo " kubectl not found, please install:" \
	echo " for OS X: brew install kubectl" \
	echo " for Ubuntu: sudo snap install kubectl --classic"
        exit 0
 fi

echo "=> everything ready, let's start"
printf '\n'
# setup vm and install microk8s
echo " setup VM and install microk8s into ..." 
printf '\n'
multipass launch --name microk8s-vm --mem 4G --disk 60G -c 4 &&
sleep 10

multipass exec microk8s-vm -- sudo snap install microk8s --classic 
sleep 5

multipass exec microk8s-vm -- sudo iptables -P FORWARD ACCEPT
sleep 15

multipass exec microk8s-vm -- /snap/bin/microk8s.enable dns dashboard
sleep 5

multipass exec microk8s-vm -- /snap/bin/microk8s.enable storage
sleep 10 

multipass exec microk8s-vm -- /snap/bin/microk8s.enable dns ingress
sleep 3

multipass exec microk8s-vm -- /snap/bin/microk8s.config > ~/kubeconfig 
export KUBECONFIG=$KUBECONFIG:~/kubeconfig

# setup kubectl
if ! grep -q KUBECONFIG "~/.bashrc"; then
 	cat 'export KUBECONFIG=$KUBECONFIG:~/kubeconfig' >> ~/.bashrc && . ~/.bashrc
     else
  echo " KUBECONFIG set, ignoring ..."
fi
printf '\n'

multipass list &&
printf '\n'

# check if we can reach kubernetes and abort if not
 if kubectl cluster-info; then
	echo "kubernetes running ..."
    else
	echo "something went wrong, check the logs, aborting "
	exit 0
 fi
	 	
echo " installing infinimesh operator "
printf '\n'
kubectl apply -f https://raw.githubusercontent.com/infinimesh/operator/master/manifests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/infinimesh/operator/master/manifests/operator.yaml
sleep 2

echo " installing kubeDB from https://github.com/kubedb "
printf '\n'
curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.11.0/hack/deploy/kubedb.sh | bash
sleep 5

echo " installing Kafka from https://github.com/strimzi "
printf '\n'
kubectl create namespace kafka &&
curl -L https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.11.1/strimzi-cluster-operator-0.11.1.yaml \
  | sed 's/namespace: .*/namespace: kafka/' \
  | kubectl -n kafka apply -f -
printf '\n'
echo "=> now we install infinimesh ..."
printf '\n'

kubectl apply -f https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/infinimesh-platform.yaml
kubectl apply -f https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/infinimesh-kafka.yaml -n kafka

echo " creating self - signed certificates "
printf '\n'
mkdir -p certs && cd certs
openssl genrsa -out ca.key 4096
openssl req -subj '/CN=infinimesh.local/O=Infinimesh' -new -x509 -sha256 -key ca.key -out ca.crt -days 3650
openssl genrsa -out apiserver_grpc.key 4096
openssl req -subj /CN=grpc.api.infinimesh.local -out apiserver_grpc.csr -key apiserver_grpc.key -new
openssl x509 -req -days 3650 -in apiserver_grpc.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out apiserver_grpc.crt -sha256
openssl genrsa -out apiserver_rest.key 4096
openssl req -subj /CN=api.infinimesh.local -out apiserver_rest.csr -key apiserver_rest.key -new
openssl x509 -req -days 3650 -in apiserver_rest.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out apiserver_rest.crt -sha256 
openssl genrsa -out mqtt_bridge.key 4096
openssl req -subj /CN=mqtt.api.infinimesh.local -out mqtt_bridge.csr -key mqtt_bridge.key -new
openssl x509 -req -days 3650 -in mqtt_bridge.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out mqtt_bridge.crt -sha256 
openssl genrsa -out app.key 4096
openssl req -subj /CN=app.infinimesh.local -out app.csr -key app.key -new
openssl x509 -req -days 3650 -in app.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out app.crt -sha256
sleep 3

echo " install the certificates"
printf '\n'
kubectl create secret tls apiserver-grpc-tls --cert apiserver_grpc.crt --key apiserver_grpc.key 
kubectl create secret tls apiserver-rest-tls --cert apiserver_rest.crt --key apiserver_rest.key 
kubectl create secret tls mqtt-bridge-tls --cert mqtt_bridge.crt --key mqtt_bridge.key 
kubectl create secret tls app-tls --cert app.crt --key app.key 

# getting IP and add hosts entries
echo " checking for host entries"
printf '\n'

IP=`multipass list|grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}"`
if ! grep -q infinimesh.local "/etc/hosts"; then
echo "=> please add this host entries into /etc/hosts: "
printf '\n'
cat <<EOL
# infinimesh-local
$IP grpc.api.infinimesh.local
$IP api.infinimesh.local
$IP app.infinimesh.local
$IP mqtt.api.infinimesh.local
EOL
   else
        echo "=> host entries already present"
fi
printf '\n'

echo "=> installing inf (infinimesh CLI) and point to the local setup:"
curl -L https://bit.ly/2CNKWzJ | BINDIR=$HOME/bin bash  
echo 'export PATH=$HOME/bin:$PATH' >> ~/.profile && . ~/.profile  
~/bin/inf config set-context local --apiserver grpc.api.infinimesh.local:443 --tls=true --ca-file ~/infinimesh-local/certs/ca.crt

printf '\n'
echo "we wait 30 secs to get all things in place ...."
secs=$((60))
while [ $secs -gt 0 ]; do
   echo -ne "$secs\033[0K\r"
   sleep 1
   : $((secs--))
done
printf '\n'

# export kubeconfig again, to be sure
export KUBECONFIG=$KUBECONFIG:~/kubeconfig

echo "infinimesh is now ready, point your browser to app.infinimesh.local or use our CLI"
echo "your master user credentials are: "
echo " root" 
kubectl get secret my-infinimesh-root-account -o=jsonpath='{.data.password}' | base64 -D
printf '\n'
echo "To trust the root certificate, you must go to your browser settings and add the file ca.crt as an certificate Authority. This works best with Firefox or Safari, we encounter some issues with Chrome."

echo "please read our documention under https://infinimesh.github.io/infinimesh/docs/#/ to proceed with adding devices and sending data. Happy IoTing ..."
printf '\n'
echo "starting infinimesh CLI:"
~/bin/inf login
