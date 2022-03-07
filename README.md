[![main](https://github.com/nayyara-cropsey/jwt-mock/workflows/Build/badge.svg)](https://github.com/nayyara-cropsey/jwt-mock/actions?query=workflow%3ABuild)
[![docker](https://github.com/nayyara-cropsey/jwt-mock/workflows/Docker/badge.svg)](https://github.com/nayyara-cropsey/jwt-mock/actions?query=workflow%3ADocker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

## Go Tests

The `jwtmocktest` package provides a HTTP test server similar to the `httptest` package. It can be used in tests to start
a test JWT mock HTTP server and used as the authorization server for a microservice using JWT authorization.

```go 
import (
  "github.com/nayyara-cropsey/jwt-mock/jwtmocktest"
  "github.com/nayyara-cropsey/jwt-mock/jwt"
)

// setup server
server, err := jwtmocktest.NewServer()

// wire to as an authorization server to a dependent microservice
setAuthorizationServerURL(server.URL)

// generate a JWT for use in Authorization header
token, err := server.GenerateJWT(jwt.Claims{
  "sub": "test-user",  // subject
  "iat": 1646451994 // issued-at epoch time
  "exp": 1646451994 // expiration epoch time 
})

// shutdown server 
server.Close()
```

## Docker

This image is pushed to `nayyarasamuel7/jwt-mock` repository. Follow these steps to get it running:

```bash 
docker pull nayyarasamuel7/jwt-mock:latest
docker run -p 80:80 nayyarasamuel7/jwt-mock:latest
```

You can override config values by setting environment variables. For example:

```bash 
docker run -p 80:80 --env JWT_MOCK_KEY_LENGTH=2048 nayyarasamuel7/jwt-mock:latest
```

The Docker repo also have various other immutable release tags pushed that match Git tags on this repo with the
format `v*`.
