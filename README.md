# JWT Mock

JWT Mock is a web server that is used to mock out a JWT-based authorization service used by your API/service. Mocking
out a real auth service is especially helpful in end-to-end testing so that JWTs can be signed and vended by the mock
JWT server.

These JWTs can help you cover various testing scenarios.

* Create tokens for users in your tests using with fake identities.
* Create tokens with various scopes/permissions for testing access control.

## Quick Start

Follow these 3 simple steps to use JWT Mock in your tests:

1) Start JWT mock - this will generate an RSA-256 signing key and associated JWKS for signing JWTs.
2) Configure your application to use JWT mock server. The JWKS path will be at `.well-known/jwks.json`.
3) During tests, generate JWTs using the `POST /jwt-generate` endpoint from the running mock server.

## Config

The default config for this service is as follows:

```yaml
port: 80
key_length: 1024
cert_life_days: 1 
```

You can override any of these through environment variables using the prefix `JWT_MOCK` . For example override port to
90 by using:

```bash 
export JWT_MOCK_PORT=90
```

## API Documentation

The JWT Mock API is documented using Open API/Swagger [here](./docs/oas.yaml). Be sure to reference it for any questions
about the API endpoints.

You will also find helpful examples using `curl` [here](./docs/curl_example.md).

## Docker 

