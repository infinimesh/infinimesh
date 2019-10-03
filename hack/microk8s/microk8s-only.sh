#!/bin/sh
clear

# install microk8s
sudo snap install microk8s --classic --channel=1.13/stable
sleep 30

# set sudo for accessing k8s > 1.14
sudo usermod -a -G microk8s multipass

sudo iptables -P FORWARD ACCEPT
sleep 15

/snap/bin/microk8s.enable dns dashboard ingress registry
sleep 5

/snap/bin/microk8s.enable storage
sleep 10

/snap/bin/microk8s.config > ~/kubeconfig 
export KUBECONFIG=$KUBECONFIG:~/kubeconfig

# retrieve token
token=$(kubectl -n kube-system get secret | grep default-token | cut -d " " -f1)
echo `kubectl -n kube-system describe secret $token` > .k8stoken
kubectl -n kube-system describe secret $token

# setup kubectl
if ! grep -q KUBECONFIG ~/.bashrc; then
 	echo "export KUBECONFIG=$KUBECONFIG:~/kubeconfig" >> ~/.bashrc
     else
  echo " KUBECONFIG set, ignoring ..."
fi
printf '\n'

# check if we can reach kubernetes and abort if not
 if kubectl cluster-info; then
	echo "kubernetes running ..."
    else
	echo "something went wrong, check the logs, aborting "
	exit 0
 fi

# add standard storage class for postgres op
echo " enable standard storage class and patch to non-default"
kubectl apply -f https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/storage.yaml
sleep 15

kubectl patch storageclass standard -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'
sleep 10
