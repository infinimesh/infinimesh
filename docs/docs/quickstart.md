# Quick start
## Install CLI
Infinimesh features a CLI to interact with Infinimesh API.
Install it automatically with our install script:
```
curl https://raw.githubusercontent.com/infinimesh/infinimesh/master/godownloader.sh | BINDIR=$HOME/bin bash
```
Please note: Depending on your OS/distribution, it may be necessary to add ~/bin to your PATH.
```
export PATH=$HOME/bin:$PATH
```
Add this to your `~/.profile` or `~/.zshrc`, depending on your shell.
After installing, set up the CLI to use our managed SaaS offering by running this command:
```
inf config set-context saas --apiserver grpc.api.infinimesh.io:443 --tls=true
```
This adds a ```context``` entry to the CLI configuration at `~/.inf/config`. TLS is enforced to guarantee secure communication with the API server.

Now, you can log into infinimesh. Run ```inf login``` and enter your username and password. This requests a token from the API server - your username and password is *NOT* stored - we take security very seriously.

## Creating a device
To get started, we will create a device. Please note, every device _needs_ to have an own certificate.

At present, Infinimesh supports only X509 certificate authentication for devices.
Generate a private key for the device:
```
openssl genrsa -out sample_1.key 4096
```
Generate the client certificate (and self-sign it):
```
openssl req -new -x509 -sha256 -key sample_1.key -out sample_1.crt -days 365
```
A namespace for your user will be automatically created, the name of the namespace is your username. To show the namespace(s) use the list command:
```
inf namespace list
```
Register the device in infinimesh's device registry:
```
inf device create my-sample-device --cert-file sample_1.crt
```
The device is registered and the fingerprint of the certificate is returned. The platform uses this fingerprint to uniquely identify your device. To check your devices and get the UID use the list command:
```
inf device list
ID     NAME               ENABLED
0x9c   my-sample-device   value:true
```
## Send states from a device to infinimesh
To simulate a device, we use the mosquitto_pub client. You can use any MQTT client, e.g. eclipse paho as well as Microsoft Edge on RaspberryPI, Yocto MQTT layers or Ubuntu Core based snaps. We use sometimes MQTTBox (http://workswithweb.com/html/mqttbox/installing_apps.html).
```
mosquitto_pub --cafile /etc/ssl/certs/ca-certificates.crt --cert sample_1.crt --key sample_1.key -m '{"abc" : 1337}' -t "devices/<YOUR DEVICE ID>/state/reported/delta" -h mqtt.api.infinimesh.io  --tls-version tlsv1.2 -d -p 8883
```
Please note: you will most likely need to adjust the --cafile flag to the path to the CA certificates on your system. This is OS specific.

You will see the following output:
```
Client mosqpub|3396-thinkpad sending CONNECT
Client mosqpub|3396-thinkpad received CONNACK (0)
Client mosqpub|3396-thinkpad sending PUBLISH (d0, q0, r0, m1, 'devices/0x1/state/reported/delta', ... (14 bytes))
Client mosqpub|3396-thinkpad sending DISCONNECT
```
The data has been sent successfully to the platform. To send more than one value per API call you can use JSON arrays in any complexity (https://www.w3schools.com/js/js_json_arrays.asp).

## Read device data from the platform
We managed to send data from a device to the platform. Now let's read back the device data from infinimesh!
You can do this via gRPC or HTTP API. The simplest way is with the CLI (which uses gRPC).
```
inf state get 0x8a
```
Replace `0x8a` with the ID of your device.
The output will look like this:
```
Reported State:
  Version:    2
  Timestamp:  2019-03-30 20:53:41.844158131 +0100 CET
  Data:
    {
      "abc": 1337
    }
Desired State: <none>
Configuration: <none>
```
In this case, the device sent a datapoint `abc` with the value `1337` at `20:53:41`.

## Send states from the platform to the device
Sending states (`desired states`) to a device is very simple. You only need to know the deviceID.
Use the API, or just the CLI.

```
inf state set 0x9c 1337
```

This sends the state `1337` to the device. Note that repeatedly sending the same state does not trigger a new message every time. Only changes are sent to the device.

Once the state has been sent, you can inspect it on the server by running:
```
inf state get 0x9c
```

The state is visible in the `desired` section.

## Thank you

Many thanks for your time and if you have any questions don't hesitate to get in touch with us! We are grateful for any improvements to the platform or this documentation, just send us a PR.
