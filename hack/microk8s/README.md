# microk8s automated setup
This script installs infinimesh and microk8s to allow a seamless developer experience. We support at the moment OS X (well tested) and Ubuntu (not so well tested). We need two tools installed on your system:
```
multipass
kubectl
```
## Quickstart
```
bash <(curl -s https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/infinimesh-setup.sh)
```

The script installs microk8s in his own VM, enables DNS, Storage, Ingress + Grafana monitoring and infinimesh as a local development system. Please read our documention https://infinimesh.github.io/infinimesh/docs/#/ and our blog https://blog.infinimesh.io/ for latest insights and HowTos.

Thank you for using infinimesh, you rock!
