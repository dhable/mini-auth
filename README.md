# mini-auth

A simple JWT/JWKS service implementation in go as a learning exercise.

## Background

This is a simple service with a sqlite DB designed to show off various aspects
of go and how a JWT/JWKS service might work. __THIS SHOULD NOT BE USED IN PRODUCTION!__
Please do not use this as a basis for a production auth service. Instead find a well
supported, open source package that is actively maintained or use a dedicated auth
provider, like Auth0.

The service uses a sqlite database for storage and is seeded with a series of DB migrations
in the `migrations` directory. Additional data or schema changes should be made through
new migrations.

The code can support roatating keys in the `jwt.KeyStore` type. The code to roate the keys
is currently missing and so it only holds a single key that does not rotate.

## API

### POST /authenticate

Given an email address and password, validates if the password for a user matches the bcrypt
hashed value stored in the sqlite database. Returns HTTP/401 Unauthorized if the values do
not match.

Request Body:
```
{
   "email": "...",
   "password": "..."
}
```

Response Body:
```
{
    "jwt": "..."
}
```

### GET /profile

Returns the profile fields for the user identified by the JWT token in the `Authorization` header.
Returns HTTP/401 Unauthorized if the header is missing, the JWT is expired, incorrect or if the
user profile doesn't exist. 

Response Body:
```
{
    "email": "...",
    "name": "...",
    "location": "..."
}
```

### GET /jwks

Returns the active set of public keys used to sign JWT tokens in the JWKS format. See RFC 7517
for response format details.


## Running the Service

You can start the service with `go run .` from the project root directory.
