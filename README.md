# Alpine
Is a microservice which provides credential (pin) service.
Credential are stored at hashed format with salt.

![Build](https://github.com/pauluswi/alpine/actions/workflows/build.yml/badge.svg)
[![codecov](https://codecov.io/gh/pauluswi/alpine/branch/master/graph/badge.svg)](https://codecov.io/gh/pauluswi/alpine)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


---

Hashed Pin = base64 encoded[sha512(pin+ salt)]

Example:

Origin Pin = 123456

Hashed Pin = ddivu2QSb2hIk9uGpihW3v2KwkEoHeyIk7GiE96OfxeuXXIK9JazAeWmvbtQLCiZJVLz7qWZoQUN9vJ3WJHOHw==

---
A salt is added to the hashing process to force their uniqueness, increase their complexity without increasing user requirements, and to mitigate password attacks like hash tables


## Project Layout

Alpine uses the following project layout:
 
```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── pin              pin-related features
│   ├── auth             authentication feature
│   ├── config           configuration library
│   ├── entity           entity definitions and domain logic
│   ├── errors           error types and handling
│   ├── healthcheck      healthcheck feature
│   └── test             helpers for testing purpose
├── migrations           database migrations
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── graceful         graceful shutdown of HTTP server
│   ├── log              structured and context-aware logger
│   └── pagination       paginated list
└── testdata             test data scripts
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example, 
the `pin` directory contains the application logic related with the payment token feature. 

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

# Getting Started

```shell
# download the repo
git clone https://github.com/pauluswi/alpine.git

cd alpine

# start a PostgreSQL database server in a Docker container
make db-start

# seed the database with some test data
make testdata

# run the RESTful API server
make run
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. 
It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /v1/login`: authenticates a user and generates a JWT
* `POST /v1/create`: create a pin
* `POST /v1/change`: change a pin
* `POST /v1/validate`: validate the pin is valid or not
* `GET /v1/get/:customer_id`:return pin information belongs to a customer

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.

If you have `cURL` or some API client tools (e.g. [Postman](https://www.getpostman.com/)), you may try the following 
more complex scenarios:

```shell
# authenticate the user via: POST /v1/login
curl -X POST -H "Content-Type: application/json" -d '{"username": "demo", "password": "pass"}' http://localhost:8080/v1/login
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, access the album resources, such as: GET /v1/xxx
# start example
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/xxx
# end example 

```

## Updating Database Schema

The starter kit uses [database migration](https://en.wikipedia.org/wiki/Schema_migration) to manage the changes of the 
database schema over the whole project development phase. The following commands are commonly used with regard to database
schema changes:

```shell
# Execute new migrations made by you or other team members.
# Usually you should run this command each time after you pull new code from the code repo. 
make migrate

# Create a new database migration.
# In the generated `migrations/*.up.sql` file, write the SQL statements that implement the schema changes.
# In the `*.down.sql` file, write the SQL statements that revert the schema changes.
make migrate-new

# Revert the last database migration.
# This is often used when a migration has some issues and needs to be reverted.
make migrate-down

# Clean up the database and rerun the migrations from the very beginning.
# Note that this command will first erase all data and tables in the database, and then
# run all migrations. 
make migrate-reset
```


### Managing Configurations

The application configuration is represented in `internal/config/config.go`. When the application starts,
it loads the configuration from a configuration file as well as environment variables. The path to the configuration 
file is specified via the `-config` command line argument which defaults to `./config/local.yml`. Configurations
specified in environment variables should be named with the `APP_` prefix and in upper case. When a configuration
is specified in both a configuration file and an environment variable, the latter takes precedence. 

The `config` directory contains the configuration files named after different environments. For example,
`config/local.yml` corresponds to the local development environment and is used when running the application 
via `make run`.

Do not keep secrets in the configuration files. Provide them via environment variables instead. For example,
you should provide `Config.DSN` using the `APP_DSN` environment variable. Secrets can be populated from a secret
storage (e.g. HashiCorp Vault) into environment variables in a bootstrap script (e.g. `cmd/server/entryscript.sh`)

## Unit Testing and Its Coverage

For testability purpose, unit testings are provided.
We can use golang test package.

```shell
$ go test -v internal/pin/*.go -race -coverprofile=coverage.out

```

## Deployment

The application can be run as a docker container. You can use `make build-docker` to build the application 
into a docker image. The docker container starts with the `cmd/server/entryscript.sh` script which reads 
the `APP_ENV` environment variable to determine which configuration file to use. For example,
if `APP_ENV` is `qa`, the application will be started with the `config/qa.yml` configuration file.

You can also run `make build` to build an executable binary named `server`. Then start the API server using the following
command,

```shell
./server -config=./config/prod.yml
```

## Reference

Go RESTful API Starter Kit (Boilerplate)
https://github.com/qiangxue/go-rest-api




