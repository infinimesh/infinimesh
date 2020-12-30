# Device State Management

Here you can manage particular device.

## State Card

you can see the Device State Card after selecting your device on UI.
It has two columns: **Reported** and **Desired** state:

![Device State](../UI/Images/device/state-base.jpg?raw=true)

**Mark 1 - Reported** state is the state received from the device.

Here you can see a last report timestamp and "version" - order number(**Marks 3 and 2**)

Device state can be reported using Eclipse **mosquitto_pub**.

**MQTT version 3.1.1/3.1 client**

Example : mosquitto_pub --cafile cert.pem --cert test.crt --key test.key  -t “devices/{device_id}/state/reported/delta" -h mqtt.api.infinimesh.cloud  --tls-version tlsv1.2 -V mqttv5 -d -p 8883 -m "{\"ping\": \"test\”}"

**MQTT version 5 client**

Example : mosquitto_pub --cafile cert.pem --cert test.crt --key test.key  -t “devices/{device_id}/state/reported/delta" -h mqtt.api.infinimesh.cloud  --tls-version tlsv1.2 -V mqttv5 -d -p 8883 -m "{"Timestamp":"","Message":[{"Topic":"T0","Data":{"ping":"test"}}]}"

By clicking on **Edit** button(**Mark 7**) - you enter **Desired** state edit mode(JSON editor - **Mark 1** below) - this is the data to be sent to the device.

![Device State Edit Mode](../UI/Images/device/state-edit-mode.jpg?raw=true)

**Desired** State can be subscribed using Eclipse **mosquitto_sub**

mosquitto_sub --cafile cert.pem --cert test.crt \
         --key test.key  -t "devices/{device_id}/state/desired/full" -h mqtt.api.infinimesh.cloud  --tls-version tlsv1.2 -V mqttv311 -d -p 8883

Note : Information on cafile, cert and key creation can be found under [Device Certificate Creation](Technical/CertificateCreation.md)

