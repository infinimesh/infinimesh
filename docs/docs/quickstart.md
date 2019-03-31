# Quick start
## Install CLI
Infinimesh features a CLI to interact with our API.
Install it automatically with our install script:
```
curl https://raw.githubusercontent.com/infinimesh/infinimesh/master/godownloader.sh | BINDIR=$HOME/bin bash
```
Please note Depending on your OS/distribution, it may be necessary to add ~/bin to your PATH.
```
export PATH=$HOME/bin:$PATH
```
Add this to your `~/.profile` or `~/.zshrc`, depending on your shell.
After installing, set up the CLI to use our managed SaaS offering by running this command:
```
inf config set-context saas --apiserver grpc.api.infinimesh.io:443 --tls=true
```
This adds a ```context``` entry to the CLI configuration at `~/.inf/config`. TLS is enforced to guarantee secure communication with the API server.

Now, you can log in to infinimesh. Run ```inf login``` and enter your username and password. This requests a token from the API server - your username and password is *NOT* stored - we take security very seriously.

## Creating a device
To get started, we'll create a device. Please note, every device _needs_ to have a own certificate. We strongly advice to use human readable names like raspi-building-campus1 for key and certificate generation. It makes the later work much easier. In future we will implement bulk device creation for ODM factory deployments.
 
At this time, Infinimesh supports only X509 certificate authentication for devices. 
Generate a private key for the device:
```
openssl genrsa -out sample_1.key 4096
```
Generate the client certificate (and self-sign it):
```
openssl req -new -x509 -sha256 -key sample_1.key -out sample_1.crt -days 365
```
Register the device in infinimesh's device registry:
```
inf device create my-sample-device --cert-file sample_1.crt -n <YOUR USERNAME>
```
The device is registered and the fingerprint of the certificate is returned. The platform uses this fingerprint to uniquely identify your device. To check your devices and get the UID use the list command:
```
inf device list
ID     NAME               ENABLED
0x9c   my-sample-device   value:true
```
## Send data and receive data from infinimesh
To simulate a device, we use the mosquitto_pub client. You can use any MQTT client, e.g. eclipse paho as well as Microsoft Edge on RaspberryPI, Yocto MQTT layers or Ubuntu Core based snaps. We use sometimes MQTTBox (http://workswithweb.com/html/mqttbox/installing_apps.html).
```
mosquitto_pub --cafile /etc/ssl/certs/ca-certificates.crt --cert sample_1.crt --key sample_1.key -m '{"sample-datapoint" : 1337}' -t "shadows/<YOUR DEVICE ID>" -h mqtt.api.infinimesh.io  --tls-version tlsv1.2 -d -p 8883
```
Please note: you will most likely need to adjust the --cafile flag to the path to the CA certificates on your system. This is OS specific.

You will see the following output:
```
Client mosqpub|3396-thinkpad sending CONNECT
Client mosqpub|3396-thinkpad received CONNACK (0)
Client mosqpub|3396-thinkpad sending PUBLISH (d0, q0, r0, m1, 'shadows/0x1', ... (14 bytes))
Client mosqpub|3396-thinkpad sending DISCONNECT
```
The data has been sent successfully to the platform. To send more as one value per API call you can use JSON arrays in any complexity (https://www.w3schools.com/js/js_json_arrays.asp). Here is an JSON example from one of our BACnet tests:
```
{
"temperature":72,
"t_metric":"F",
"co2":1456,
"c_metric":"ppm",
"noise":56,
"n_metric":"db",
"spot_enabled":3,
"spot_light_brightness":[ "78", "73", "44" ],
"s_metric":"lux"
}
```
 
## Read device data from the platform
We managed to send data from a device to the platform. Now let's read back the device data from infinimesh!
You can do this via gRPC or HTTP API. The simplest way is with the CLI.
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
In this case, the device sent a datapoint `abc` with the value `1337` at `20:53:41`. Remember the BACnet data? That looks like this:
```
inf state get 0x9c
Reported State:
  Version:    5
  Timestamp:  2019-03-31 12:47:13.158610661 +0200 DST
  Data:
    {
      "c_metric": "ppm",
      "cmetric": "ppm",
      "co2": 1456,
      "metric": "ppm",
      "n_metric": "db",
      "noise": 56,
      "s_metric": "lux",
      "spot_enabled": 3,
      "spot_light_brightness": [
        "78",
        "73",
        "44"
      ],
      "t_metric": "F",
      "temperature": 72,
      "tmetric": "F"
    }
Desired State: <none>
Configuration: <none>
```

Thank you for your time and if you have any questions don't hesitate to get in touch with us!

