# Generate Authentication Token

All the REST Endpoints for the product requries an Authentication token for them to work. Otherwise you will get below error if the token is not provided.

> Request unauthenticated with bearer

## How to obtain the Token 

Pre-Requisites: 

1. You need valid user credentials for the applications

Steps:

1. REST Request Details for Token generation
   
   - REST Endpoint: **{URL}/account/token**
   > URL is the domain for the environment E.g. console.infinimesh.dummy
   - Request Type: **POST**
   - Request Header: **Content-Type: application/json**
   - Request Body: **Content in JSON format**

2. Once the above REST Request is send with the required JSON body to the endpoint, a token will be generated and sent back in JSON format.

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
   - Request Header: **Authorization: bearer Authentication_Token**


