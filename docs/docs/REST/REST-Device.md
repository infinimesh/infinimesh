# Device Registry Endpoint

The Device Registry Endpoint allows you to mange devices for the applications. Below are the endpoints avaiable:

| Endpoints | Purpose of the Endpoint |
|-----------|-------------------------|
| /devices?namespaceid=<namespaceid> | Get details of all devices |
| /devices | Create a device |
| /devices/{device.id} | Update a device |
| /devices/{deviceid}/owner/{ownerid} | Add an owner to the device |
| /devices/{deviceid}/owner/{ownerid} | Remove an owner from the device |
| /devices/{id} | Get details of a specific device |
| /devices/{id} | Delete a specific device |

> Request unauthenticated with bearer

## How to obtain the Token 

Pre-Requisites: 

1. You need valid user credentials for the applications

Steps:

1. REST Request Details for Token generation
   
   - REST Endpoint: **<URL>/account/token**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **POST**
   - Request Header: **Content-Type: application/json**
   - Request Body: **Content in JSON format**

2. Once the above REST Request is send to the endpoint, a token will be generated and sent back in JSON format.

Example:

Sample Request Body:
```
{
"password": "Enter Password here",
"username": "Enter UserID here"
}
```

Sample Response:
```
{
"token": "Authentiction_Token"
}
```

3. The token received in response will be used for all the other REST request. The token has to be set in the REST Request header. Below are the details:
   - Request Header: **Authorization: bearer Authentiction_Token**


