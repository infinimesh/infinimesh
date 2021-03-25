# Device States Endpoint

The Device States Endpoints allows you to manage the states of devices. Below are the endpoints avaiable:

| HTTP Request | Endpoints | Purpose of the Endpoint |
|--------------|-----------|-------------------------|
| GET | /devices/{id}/state
| PATCH | /devices/{device.id}/state
| GET | /devices/{id}/state/stream

## How to GET the Device State

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a device id for which you want to see the current state

Steps:

1. REST Request Details to get a Device state
   
   - REST Endpoint: **<URL>/devices/{id}/state
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **id should be a valid device id**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the correct device id and correct token, the specific device reported state will be returned back with HTTP 200 response. Otherwise you will get an error with the reason why the device reported state was not returned.

## How to PATCH the desired Device State

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a device id to which you want to patch the desired state

Steps:

1. REST Request Details to patch the desired Device state

   - REST Endpoint: **<URL>/devices/{id}/state
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **id should be a valid device id**
   - Request Type: **PATCH**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the correct device id and correct token, the specific desired device state will be patched and HTTP 200 response will be recieved. Otherwise you will get an error with the reason why the desired state patching was not successful.

## How to get the device reported state streaming

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a device id to which you want to stream the device reported state

Steps:

1. REST Request Details to get the device reported state streaming

   - REST Endpoint: **<URL>/devices/{id}/state/stream
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **id should be a valid device id**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the correct device id and correct token, the specific device reported state will be streamed and HTTP 200 response will be recieved. Otherwise you will get an error with the reason why the streaming was not successful.

curl request for streaming : curl -X GET "https://<URL>/devices/<device id>/state/stream?only_delta=false" -H "Authorization: Bearer <Authentication_Token>"