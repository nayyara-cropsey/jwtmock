# Examples

Examples using CURL and running JWT Mock server locally on port 80.

## Create JWT Token

```
curl --request POST \
  --url http://localhost/jwtmock/generate-jwt \
  --header 'Content-Type: application/json' \
  --data '{
  "iss": "https://auth.mine.go/",
  "sub": "nayyara",
  "aud": [
    "https://api.mine.go"
  ],
  "iat": 1643651018,
  "exp": 1643748940  
}
'
```

JWTs can be decoded at [jwt.io](https://jwt.io).

## Get JWKS

```
curl --request GET \
  --url http://localhost/.well-known/jwks.json
```

## Refresh JWKS

```
curl --request POST \
  --url http://localhost/.well-known/jwks.json
```
