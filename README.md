# skills-api

## What is this?
This is a practice project that I setup to go back and confirm my knolledge of Golang by creating a basic API that uses Docker, PostgreSQL (with an ORM), and uses test driven development.

This is very early stage and is not complete. Working on establishing a good pattern for TDD, then I need to explore the Auth options. After that I can finish up the basic API entities. I may also go back after I have the REST/CRUD endpoints done and figure out if I can also apply a GraphQL endpoint to this example.

### Technical Details
This is a an API that us written in Go Language. It uses the Go package `Fiber` as the API framework. It uses the package `Grom` as the ORM. It uses a PostgresSQL database. It uses Docker container, with a seperate database server and web server. This also uses the package `Air` for hot relaods in the web server.

### Purpose
The idea of this API is that you can use this as a data resource to store team information of skills that a team has and can be linked to projects.

## Running the Project
To start this project, you should be abel to downlaod the source, then copy the `.env.template` file to a new `.env` file (use `cp .env.template .env`). Then change the values in the `.env` fiel to match what your preferences. This should be at least the database information.

Once the `.env` file is setup, you should be able to use the make command to start the Docker container.
```
make dev
```

## Test Driven Development (TDD)
Currently you can run the test from teh `api` folder by using the command `make test`. This should run all the current tests for the API endpoints.
When you make changes to the project, it is suggeste dthat you start by creating a new test file in the `api/tests` folder. The tests should use an in memory Sqlite database that uses the same schema of the PostgreSQL database.