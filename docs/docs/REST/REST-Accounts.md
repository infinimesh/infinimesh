# Accounts Endpoint

The Accounts Endpoint allows you to mange user accounts for the applications. Below are the endpoints avaiable:

| HTTP Request | Endpoints | Purpose of the Endpoint |
|--------------|-----------|-------------------------|
| POST | /accounts | Create an user account |
| GET | /accounts/{id} | Get details of a specific account |
| PATCH | /accounts/{account.uid} | Update an account |
| DELETE | /accounts/{uid} | Delete a specific account |
| GET | /accounts | Get details of all user accounts |
| GET | /account | Get details of current user |
| PUT | /accounts/{accountid}/owner/{ownerid} | Add an owner to the user account |
| DELETE | /accounts/{accountid}/owner/{ownerid} | Remove an owner from the user account |



## How to create an Account

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)

Steps:

1. REST Request Details for Creating an Account
   
   - REST Endpoint: **<URL>/accounts**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **POST**
   - Request Header: **Content-Type: application/json**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "account": {
    "default_namespace": {
      "deleteinitiationtime": "string",
      "id": "string",
      "markfordeletion": false,
      "name": "string"
    },
    "enabled": false,
    "is_admin": false,
    "is_root": false,
    "name": "string",
    "owner": "string",
    "password": "string",
    "uid": "string",
    "username": "string"
  },
  "create_gf_user": false,
  "password": "string"
}
```

Sample Body Format:
```
{
  "account": {
    "enabled": true,
    "is_admin": false,
    "name": "FirstName LastName",
    "password": "NewPassword",
    "username": "NewUser"
  }
}
```

2. Once the above REST Request is send with the required JSON body to the endpoint, an account will be created and JSON response will be send back.

Sample Response:
```
{
  "uid":"string"
}
```

## How to get an Account data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace

Steps:

1. REST Request Details for getting a namespace data
   
   - REST Endpoint: **<URL>/accounts/{id}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **id should be a valid account id**
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameters to the endpoint, if the account id is valid then the account details in JSON response will be send back.

Response Format:
```
{
  "account": [
  {
    "default_namespace": {
      "deleteinitiationtime": "string",
      "id": "string",
      "markfordeletion": false,
      "name": "string"
    },
    "enabled": false,
    "is_admin": false,
    "is_root": false,
    "name": "string",
    "owner": "string",
    "password": "string",
    "uid": "string",
    "username": "string"
  }
  ]
}
```

## How to update an Account

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace 

Steps:

1. REST Request Details for Updating an Account
   
   - REST Endpoint: **<URL>/accounts/{account.uid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **account.uid should be a valid account id**
   - Request Type: **PATCH**
   - Request Header: **Content-Type: application/json**
   - Request Header: **Authorization: bearer Authentication_Token**
   - Request Body: **Content in JSON format**

Request Body Format:
```
{
  "default_namespace": {
    "deleteinitiationtime": "string",
    "id": "string",
    "markfordeletion": false,
    "name": "string"
  },
  "enabled": false,
  "is_admin": false,
  "is_root": false,
  "name": "string",
  "owner": "string",
  "password": "string",
  "uid": "string",
  "username": "string"
}
```

Sample Request Body:

Below is an example of an update JSON request which will update the account with ID 0x000. The fields it will update are Name, Password and Enabled.
```
{
  "enabled": true,
  "name": "NewName",
  "password": "NewPassword",
  "uid": "0x000",
}
```

2. Once the above REST Request is send with the required JSON body to the endpoint, an HTTP 200 reposne is receive if the namespace update was successful. Otherwise you will get an error with the reason why the update was not successful.

## How to delete an Account   

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace 

Steps:

1. REST Request Details for Deleting an Account
   
   - REST Endpoint: **<URL>/accounts/{uid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **uid should be a valid account id in infinimesh**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameters to the endpoint, an HTTP 200 reposne is receive if the account is deleted. Otherwise you will get an error with the reason why the request was not successful.

## How to get all Accounts data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace 

Steps:

1. REST Request Details for Getting all Accounts
   
   - REST Endpoint: **<URL>/accounts**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send to the endpoint, it returns all the accounts present. If the user is admin then all accounts manged by the Admin will be returned otherwise only the self account will be returned.

Response Format:
```
{
  "account": [
      {
    "default_namespace": {
      "deleteinitiationtime": "string",
      "id": "string",
      "markfordeletion": false,
      "name": "string"
    },
    "enabled": false,
    "is_admin": false,
    "is_root": false,
    "name": "string",
    "owner": "string",
    "password": "string",
    "uid": "string",
    "username": "string"
    },
  ]
}
```

## How to get self or current account data

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace 

Steps:

1. REST Request Details for getting all users for a Namespace
   
   - REST Endpoint: **<URL>/account**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **GET**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send to the endpoint, the details of all the current user will be displayed for the namespace in JSON format.

Response Format:
```
{
  "account": {
    "default_namespace": {
      "deleteinitiationtime": "string",
      "id": "string",
      "markfordeletion": false,
      "name": "string"
    },
    "enabled": false,
    "is_admin": false,
    "is_root": false,
    "name": "string",
    "owner": "string",
    "password": "string",
    "uid": "string",
    "username": "string"
  }
}
```

## How to add an owner for an Account

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace 

Steps:

1. REST Request Details for adding an owner for an Account
   
   - REST Endpoint: **<URL>/accounts/{accountid}/owner/{ownerid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **accountid should be a valid account id and the ownerid should be a valid user account with Admin priviledges in infinimesh**
   - Request Type: **PUT**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameters to the endpoint, an HTTP 200 reponse is receive if the request was successful. Otherwise you will get an error with the reason why the request was not successful.

## How to remove an owner from an Account

Pre-Requisites: 

1. You need valid user credentials for the applications to obtain token (Refer [here](https://infinimesh.github.io/infinimesh/docs/#/REST/GenerateToken#how-to-obtain-the-token) on how to generate a token)
2. You need a namespace 

Steps:

1. REST Request Details for removing owner from an Account
   
   - REST Endpoint: **<URL>/accounts/{accountid}/owner/{ownerid}**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Path Parameters: **accountid should be a valid account id and the ownerid should be a valid user account with Admin priviledges in infinimesh**
   - Request Type: **DELETE**
   - Request Header: **Authorization: bearer Authentication_Token**

2. Once the above REST Request is send with the required path parameter to the endpoint, the specific owner will be removed from the account and an HTTP 200 response will be received. Otherwise you will get an error with the reason why the request was not successful.

