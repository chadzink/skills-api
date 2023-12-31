# skills-api

## What is this?
This is a practice project that I setup to go back and confirm my knowledge of Golang by creating a basic API that uses Docker, PostgreSQL (with an ORM), and uses Test Driven Development (TDD).

This is very early stage and is not complete. Working on establishing a good pattern for TDD, then I need to explore the Auth options. After that I can finish up the basic API entities. I may also go back after I have the REST/CRUD endpoints done and figure out if I can also apply a GraphQL endpoint to this example.

### Technical Details
This is an API that is written in Go Language. It uses the Go package `Fiber` as the API framework. It uses the package `Grom` as the ORM. It uses a PostgreSQL database. It uses Docker container, with a separate database server and web server. This also uses the package `Air` for hot reloads in the web server.

### Purpose
The idea of this API is that you can use this as a data resource to store team information of skills that a team has and can be linked to projects.

## Running the Project
System setup will need to have Golang version 1.21 installed. You will also need docker for the PostgreSQL setup. This coul dbe skipped by installing Postgress locally and running the project locally with an environmental variables to connect to you rlocal database server.

##### Requires
1. Go Language [v 1.21 ](https://go.dev/doc/install)
2. Docker [install](https://docs.docker.com/engine/install/)
3. Local (no docker)
   * [PostgreSQL Server](https://www.postgresql.org/download/)
   * Sqlite (for testing)

To start this project, you should be able to download the source, then copy the `.env.template` file to a new `.env` file (use `cp .env.template .env`). Then change the values in the `.env` file to match your system preferences. This should be at least the database information.

Once the `.env` file is setup, you should be able to use the make command to start the Docker container.
```
make dev
```

## Test Driven Development (TDD)
Currently you can run the test from the `api` folder by using the command `make test`. This should run all the current tests for the API endpoints.
When you make changes to the project, it is suggested that you start by creating a new test file in the `api/tests` folder. The tests should use an in memory Sqlite database that uses the same schema of the PostgreSQL database.

## API Documentation
The API documentation is done with Swagger OpenAPI specifications using [swaggo/fiber-swagger](https://github.com/swaggo/fiber-swagger).

## TO DO
Look into validation example in: https://dev.to/koddr/build-a-restful-api-on-go-fiber-postgresql-jwt-and-swagger-docs-in-isolated-docker-containers-475j