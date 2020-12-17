# Device Registry Endpoint

The Device Registry Endpoint allows you to mange devices for the applications. Below are the endpoints avaiable:

| HTTP Request | Endpoints | Purpose of the Endpoint |
|--------------|-----------|-------------------------|
| POST | /devices | Create a device |
| GET | /devices/{id} | Get details of a specific device |
| PATCH | /devices/{device.id} | Update a device |
| DELETE | /devices/{id} | Delete a specific device |
| GET | /devices?namespaceid={namespaceid} | Get details of all devices |
| PUT | /devices/{deviceid}/owner/{ownerid} | Add an owner to the device |
| DELETE | /devices/{deviceid}/owner/{ownerid} | Remove an owner from the device |



## How to create a Device

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for Creating a Device
   
   - REST Endpoint: **<URL>/devices**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **POST**
   - Request Header: **Content-Type: application/json**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
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
  "device": {
    "certificate": {
      "algorithm": "string",
      "fingerprint": "string",
      "fingerprintAlgorithm": "string",
      "pem_data": "string"
    },
    "enabled": false,
    "id":"string"
    "name": "string",
    "namespace": "string",
    "tags": [
      "string"
    ]
  }
}
```

## How to get a Device data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce with a device in it

Steps:

1. REST Request Details for getting a Device data
   
   - REST Endpoint: **<URL>/devices/{id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **id should be a valid device id**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameters to the endpoint, if the device id is valid then the device details in JSON response will be send back.

Response Format:
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
    "id":"string"
    "name": "string",
    "namespace": "string",
    "tags": [
      "string"
    ]
  }
}
```

## How to update a Device

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce with a valid device in it

Steps:

1. REST Request Details for Updating a Device
   
   - REST Endpoint: **<URL>/devices/{device.id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **device.id should be a valid device id**
   - Request Type: **PATCH**
   - Request Header: **Content-Type: application/json**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "certificate": {
    "algorithm": "string",
    "fingerprint": "string",
    "fingerprintAlgorithm": "string",
    "pem_data": "string"
  },
  "enabled": false,
  "id": "string",
  "name": "string",
  "namespace": "string",
  "tags": [
    "string"
  ]
}
```

Sample Request Body:

Below is an example of an update JSON request which will update the device with ID 0x000. The fields it will update are Name, Namespace, Status (enable field) and Tags.
```
{
  "enabled": true,
  "id": "0x000",
  "name": "New_Device_Name",
  "namespace": "0x000",
  "tags": [
    "New_Tag"
  ]
}
```

2. Once the above REST Request is send with the required JSON body to the endpoint, an HTTP 200 reposne is receive if the device update was successful. Otherwise you will get an error with the reason why the update was not successful.

## How to delete a Device

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce with a valid device in it

Steps:

1. REST Request Details for Deleting a Device
   
   - REST Endpoint: **<URL>/devices/{id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **id should be a valid device id**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameter to the endpoint, the specific device will be deleted and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the deletion was not successful.

## How to get all Devices data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce with a device in it

Steps:

1. REST Request Details for Getting all Devices
   
   - REST Endpoint: **<URL>/devices?namespaceid={namespaceid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Query String Parameters: **namespaceid should be a valid namesapce id**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required querystring parameters to the endpoint, it returns all the devices present in the specified namesapce.

Response Format:
```
{
  "devices": [
      {
        "certificate": {
          "algorithm": "string",
          "fingerprint": "string",
          "fingerprintAlgorithm": "string",
          "pem_data": "string"
        },
        "enabled": false,
        "id":"string"
        "name": "string",
        "namespace": "string",
        "tags": [
          "string"
        ]
      }
   ]
}
```

## How to update a Device's owner

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce with a valid device in it

Steps:

1. REST Request Details for Updating a Device
   
   - REST Endpoint: **<URL>/devices/{deviceid}/owner/{ownerid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **deviceid should be a valid device id and the ownerid should be a valid user account in infinimesh**
   - Request Type: **PUT**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameters to the endpoint, an HTTP 200 reponse is receive if the device update was successful. Otherwise you will get an error with the reason why the request was not successful.

## How to remove a Device's owner

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce with a valid device in it

Steps:

1. REST Request Details for Deleting a Device
   
   - REST Endpoint: **<URL>/devices/{deviceid}/owner/{ownerid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **deviceid should be a valid device id and the ownerid should be a valid user account in infinimesh**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameter to the endpoint, the specific owner will be removed from the device and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the request was not successful.

