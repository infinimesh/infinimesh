# Device Registry Endpoint

The Namespace Endpoint allows you to mange namesapces for the applications. Below are the endpoints avaiable:

| HTTP Request | Endpoints | Purpose of the Endpoint |
|--------------|-----------|-------------------------|
| POST | /namespaces | Create a namespace |
| GET | /namespaces/{namespace} | Get details of a specific namespace |
| PATCH | /namespaces/{namespace.id} | Update a namespace |
| DELETE | /namespaces/{namespaceid}/{harddelete} | Delete a specific namespace |
| GET | /namespaces | Get details of all namespaces |
| GET | /namespaces/{namespace}/permissions | Get details of all the permissions for a namespaces |
| PUT | /namespaces/{namespace}/permissions/{account_id} | Add an user to the namespace |
| DELETE | /namespaces/{namespace}/permissions/{account_id} | Remove an user from the namespace |



## How to create a Namespace

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)

Steps:

1. REST Request Details for Creating a Namespace
   
   - REST Endpoint: **<URL>/namespaces**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **POST**
   - Request Header: **Content-Type: application/json**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "name": "string"
}
```

2. Once the above REST Request is send with the required JSON body to the endpoint, a namespace will be created and JSON response will be send back.

Sample Response:
```
{
  "deleteinitiationtime":"string"
  "id::"string"
  "markfordeletion":false
  "name":"string"
}
```

## How to get a Namespace data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce

Steps:

1. REST Request Details for getting a namespace data
   
   - REST Endpoint: **<URL>/namespaces/{namespace}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespace should be a valid namesapce name**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameters to the endpoint, if the device id is valid then the device details in JSON response will be send back.

Response Format:
```
{
  "deleteinitiationtime":"string"
  "id::"string"
  "markfordeletion":false
  "name":"string"
}
```

## How to update a Namespace

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for Updating a Device
   
   - REST Endpoint: **<URL>/namespaces/{namespace.id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespace.id should be a valid namespace id**
   - Request Type: **PATCH**
   - Request Header: **Content-Type: application/json**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "deleteinitiationtime": "string",
  "id": "string",
  "markfordeletion": false,
  "name": "string"
}
```

Sample Request Body:

Below is an example of an update JSON request which will update the namesapce with ID 0x000. The fields it will update are Name and Markfordeletion.
```
{
  "id": "0x000",
  "markfordeletion": true,
  "name": "NewNames"
}
```

2. Once the above REST Request is send with the required JSON body to the endpoint, an HTTP 200 reposne is receive if the namespace update was successful. Otherwise you will get an error with the reason why the update was not successful.

## How to Soft delete a Namespace 

> Note: Soft delete will mark the namespace for deletion and will only actually delete the namesapce from the Dgraph database.

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for Deleting a Namespace
   
   - REST Endpoint: **<URL>/namespaces/{namespaceid}/{harddelete}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespaceid should be a valid namesapce id**
                              **hardelete should be a set to false for soft delete**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameter to the endpoint, the specific namesapce will be marked for deletion and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the deletion was not successful.

## How to Hard delete a Namespace 

> Note: Hard delete will actually delete all the namesapces marked for deletion from the Dgraph database. The number of retention days should be over for a namespace to be hard deleted.

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for Deleting a Namespace
   
   - REST Endpoint: **<URL>/namespaces/{namespaceid}/{harddelete}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespaceid is ignored for hard delete**
                              **hardelete should be a set to true for hard delete**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameter to the endpoint, the application with run a job that will delete all the namespaces that are marked for deletion which has passed the rentention period and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the deletion was not successful.

## How to get all Namespaces data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for Getting all Namespaces
   
   - REST Endpoint: **<URL>/namespaces**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required querystring parameters to the endpoint, it returns all the devices present in the specified namesapce.

Response Format:
```
{
  "namespaces": [
      {
          "deleteinitiationtime":"string"
          "id::"string"
          "markfordeletion":false
          "name":"string"
        }
   ]
}
```

## How to get all user permissions in the Namespace

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for getting all users for a Namespace
   
   - REST Endpoint: **<URL>/namespaces/{namespace}/permissions**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespace should be a valid namespace id**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameter to the endpoint, the details of all the users will be displayed for the namespace in JSON format.

Response Format:
```
{
  "permissions": [
      {
         "account_id":"string"
         "account_name":"string"
         "action":"NONE"
         "namespace":"string"
      }
   ]
}
```

## How to add an user in the Namespace

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for adding a user to a Namespace
   
   - REST Endpoint: **<URL>/namespaces/{namespace}/permissions/{account_id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespace should be a valid namespace id and the account_id should be a valid user account in infinimesh**
   - Request Type: **PUT**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "action": "NONE"
}
```

> Note: Possible values for action field are "NONE", "READ" and "WRITE".

2. Once the above REST Request is send with the required path parameter to the endpoint, the specific user will be added to the namespace and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the request was not successful.

## How to remove an user from a Namespace

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinitedevices.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namesapce 

Steps:

1. REST Request Details for removing a user from a Namespace
   
   - REST Endpoint: **<URL>/namespaces/{namespace}/permissions/{account_id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **namespace should be a valid namespace id and the account_id should be a valid user account in infinimesh**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "action": "NONE"
}
```

> Note: Possible values for action field are "NONE", "READ" and "WRITE".

2. Once the above REST Request is send with the required path parameter to the endpoint, the specific user will be removed from the namesapce and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the request was not successful.

