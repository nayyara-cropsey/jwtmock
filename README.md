[![main](https://github.com/nayyara-cropsey/jwtmock/workflows/Build/badge.svg)](https://github.com/nayyara-cropsey/jwt-mock/actions?query=workflow%3ABuild)
[![docker](https://github.com/nayyara-cropsey/jwtmock/workflows/Docker/badge.svg)](https://github.com/nayyara-cropsey/jwt-mock/actions?query=workflow%3ADocker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# JWT Mock

JWT Mock is a web server that is used to mock out a JWT-based authorization server for an authenticated API. Mocking out
an authorization service is especially helpful in end-to-end/smoke testing to see how different handlers and middleware
respond to JWTs representing various levels of access and identities.

Mock JWTs can help you cover various testing scenarios.

* Create tokens for various user identities.
* Create tokens with various scopes/permissions for testing access control.

## Quick Start

Follow these 3 simple steps to use JWT Mock in your tests:

1) Start JWT Mock server

* Via Docker image if using a docker ecosystem
* Via `jwtmocktest.NewServer` in Go tests

2) Configure a dependent API on to use JWT Mock server as the authorization server. The server provides an endpoint to
   retrieves the JSON Web Key Set (JWKS) at `./well-known/jwks.json`

3) Generate JWTs for use in tests

* Via Docker image, use the `POST /generate-jwt` endpoint to generate a JWT from a set of claims.
* Via `jwtmocktest.Server` in Go tests, use `GenerateJWT` method on the test server.

## API Documentation

The JWT Mock API is documented using Open API and available on
SwaggerHub [here](https://app.swaggerhub.com/apis-docs/nayyara-cropsey/jwtmock/1.0.0/). Be sure to reference it for
any questions about the API endpoints.

You will also find helpful examples using `curl` [here](./docs/curl_example.md).

## Go Tests

The `jwtmocktest` package provides an HTTP test server similar to the `httptest` package. It can be used as the
authorization server for a microservice using JWTs.

```go 
import (
  "github.com/nayyara-cropsey/jwtmock/jwtmocktest"
  "github.com/nayyara-cropsey/jwtmock"
)

// setup server
server, err := jwtmocktest.NewServer()

// wire to as an authorization server to a dependent microservice's config
appConfig.AuthZServer = server.URL

// generate a JWT for use in Authorization header
token, err := server.GenerateJWT(jwtmock.Claims{
  "sub": "test-user",  // subject
  "iat": 1646451994 // issued-at epoch time
  "exp": 1646451994 // expiration epoch time 
})

// use JWT in microservice API 
req = req.WithHeader("Authorization", "Bearer: " + token)
...

// shutdown server 
server.Close()
```

Alternatively you can also use the `jwtmocktest.Client` to connect to a running JWT Mock server.

```go 
import (  
  "github.com/nayyara-cropsey/jwtmock"
)

// create client 
client, err := jwtmock.NewClient(mockJWTServerURL)

// generate a JWT for use in Authorization header
token, err := client.GenerateJWT(jwtmock.Claims{
  "sub": "test-user",  // subject
  "iat": 1646451994 // issued-at epoch time
  "exp": 1646451994 // expiration epoch time 
})

```

## Docker

This image is pushed to `nayyaracropsey/jwtmock` repository. Follow these steps to get it running:

```bash 
docker pull nayyaracropsey/jwtmock:latest
docker run -p 80:80 nayyaracropsey/jwtmock:latest
```

The Docker repo also have various other immutable release tags pushed that match Git tags on this repo with the
format `v*`.

### Config

The default config for this service is as follows:

```yaml
port: 80
key_length: 1024
cert_life_days: 1
log_level: debug
```

You can override any of these through environment variables using the prefix `JWT_MOCK`. For example override key length
using:

```bash 
docker run -p 80:80 --env JWT_MOCK_KEY_LENGTH=2048 --env nayyaracropsey/jwtmock:latest
```
