#!/bin/sh
mkdir -p ~/ibmcloud-cli && cd ~/ibmcloud-cli
# clearlinux does not have local/bin
sudo mkdir -p /usr/local/bin

# install BLuemix-CLI manually
wget -qO- https://clis.ng.bluemix.net/download/bluemix-cli/0.22.0/linux64 | tar xvz 
sh Bluemix_CLI/install

#install all tools we need
ibmcloud plugin install dev
ibmcloud plugin install cloud-object-storage
ibmcloud plugin install container-registry
ibmcloud plugin install container-service
ibmcloud plugin install cloud-functions

# disable stats data collection
ibmcloud config --usage-stats-collect false

# login
ibmcloud login

# list all clusters
ibmcloud ks cluster ls

