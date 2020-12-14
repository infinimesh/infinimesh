# Welcome to the infinimesh Backend Wiki

Here you find some notes about infinimesh REST Endpoints and how to use them. All the REST Endpoints for the product requries an Authentication token for them to work. Otherwise you will get below error if the token is not provided.

> Request unauthenticated with bearer

## How to obtain Authentication Token 

Pre-Requisites: 

1. You need valid user credentials for the applications

Steps:

1. REST Request Details for Token generation
   
   - REST Endpoint: **{{URL}}/account/token**
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
## Table of Contents

1. [Devices Registry Management](UI/REST-Device.md)
2. [Accounts Management](UI/REST-Accounts.md)
3. [Namespaces Management](UI/REST-Namespace.md)
4. [User Account Management](UI/REST-User.md)
