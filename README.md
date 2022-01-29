# JWT Mock 

JWT Mock is a web server that is used to mock out a JWT-based authorization service used by your API/service. 
Mocking out a real auth service is especially helpful in end-to-end testing so that JWTs can be signed and vended by the mock JWT server. 

These JWTs can help you cover various testing scenarios. 
* Create tokens for users in your test fixtures for consistent test data/indentities.
* Create tokens with various scopes/permissions for testing access control.

Follow these 3 simple steps to use JWT Mock in your tests:

1) Configure JWT Mock with preferred signing algorithm (key is generated automatically, you don't need to provide this)
2) Configure your application to use JWT mock server. The JWKS URL path will be `.well-known/jwks.json`.
3) During tests, generate JWTs using the `POST /jwt-generate` endpoint. 


## Generate JWT

