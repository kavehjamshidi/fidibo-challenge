# Fidibo Challenge

Golang backend service using Gin, Redis, JWT, and Docker as the technical assignment of Fidibo's hiring process.

## Prerequisites

- Go version 1.18 was used as the programming language.
- Redis is required as the data store for caching responses. Redis version 7.0.8 was used.

## Clone this project

```shell
$ git clone https://github.com/kavehjamshidi/fidibo-challenge.git

$ cd fidibo-challenge
```

## Installing Go Dependencies

All dependencies can be downloaded using the command below:

```shell
$ go mod download
```

## Setting up Environment Variables

For the sake of simplicity, Environment Variables are used for configurations. If no environment variable is provided, the fallback values which are defined as constants in the code are used.
The table below includes all required Environment Variables and their respective fallback values:
| |Environment Variable Name |Default Value |
|----------------|-------------------------------|-----------------------------|
|Redis Address|`REDIS_ADDRESS` |`localhost:6379` |
|Test Redis Address (Integration Test) |`TEST_REDIS_ADDRESS` |`localhost:6379` |
|Server Address |`SERVER_ADDRESS`|`:8080`|
|Access Token Expiry |`ACCESS_EXPIRY`|`15m`|
|Access Token Secret |`ACCESS_SECRET`|`access token secret`|
|Refresh Token Expiry |`REFRESH_EXPIRY`|`168h`|
|Refresh Token Secret |`REFRESH_SECRET`|`refresh token secret`|

## Build and Test

To run all tests, run the command below:

```shell
$ go test ./...
```

To run the project using Docker Compose, use this command:

```shell
$ docker-compose up
```

## Endpoints

_Login_ and _Refresh Token_ endpoints are public, and _Search_ endpoint is protected. To access the Search endpoint, a **JWT** token should be sent in the Authorization header. A simple **Postman** collection is included in the _docs_ directory as documentation.
For the sake of simplicity, Login endpoint always returns successful response regardless of the provided credentials. This is far from ideal and definitely not practical in real-world projects, but given the tight deadline, this was the best I could do.
Refresh Tokens are not stored in Redis or any other database. As a result, no _Logout_ functionality is present.
