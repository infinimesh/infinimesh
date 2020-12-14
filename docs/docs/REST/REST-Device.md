# Device Registry Endpoint

The Device Registry Endpoint allows you to mange devices for the applications. Below are the endpoints avaiable:

| HTTP Request | Endpoints | Purpose of the Endpoint |
|--------------|-----------|-------------------------|
| GET | /devices?namespaceid={namespaceid} | Get details of all devices |
| POST | /devices | Create a device |
| PATCH | /devices/{device.id} | Update a device |
| PUT | /devices/{deviceid}/owner/{ownerid} | Add an owner to the device |
| DELETE | /devices/{deviceid}/owner/{ownerid} | Remove an owner from the device |
| GET | /devices/{id} | Get details of a specific device |
| DELETE | /devices/{id} | Delete a specific device |

## How to create a device

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for Token generation
   
   - REST Endpoint: **<URL>/devices**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **POST**
   - Request Header: **Content-Type: application/json**
   - Request Body: **Content in JSON format**

Sample Request Body:
```
{
  "device": {
    "certificate": {
      "algorithm": "string",
      "fingerprint": "string",
      "fingerprintAlgorithm": "string",
      "pem_data": "string"
    },
    "enabled": false,
    "name": "string",
    "namespace": "string",
    "tags": [
      "string"
    ]
  }
}
```

2. Once the above REST Request is send with the required JSON body to the endpoint, a device will be created and JSON response will be send back.

Sample Response:
```
{
device:{
certificate:{
algorithm:"string"
fingerprint:"string"
fingerprintAlgorithm:"string"
pem_data:"string"
}
enabled:false
id:"string"
name:"string"
namespace:"string"
tags:[
"string"
]
}
}
```


