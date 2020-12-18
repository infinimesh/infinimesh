Here you find some notes about device certificate creation

## Certificate Creation Mode
At present, Infinimesh supports only X509 certificate authentication for devices.

Generate a **private key** for the device:

openssl genrsa -out {key_filename}.key 4096

Generate the **client certificate** (and self-sign it):

openssl req -new -x509 -sha256 -key {key_filename}.key -out {crt_filename}.crt -days 365

**Refer:** ![a link] https://github.com/infinimesh/infinimesh/blob/master/hack/device_certs/create-certs.sh

**Note:**
The generated certificate key should be same as of the one used while device creation and while sending the mosquitto message [MQTT](Technical/MQTT.md) to the device.
.


