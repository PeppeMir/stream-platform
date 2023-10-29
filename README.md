# Stream Platform application

## Prerequisites

You need to have installed the Docker engine with the ```docker compose``` plugin supporting the version 3.7 of the compose file.

## Build the app

To build the app execute

```
docker build --tag stream-platform-app .
```
After the build, check the image with name stream-platform-app and tag latest

```docker
docker image ls
```

## Startup MySQL database and Go app

To startup the whole ecosystem you just have to execute the following command:

```docker
docker compose up -d
```
The MYSQL instance will start first, by initializing some members data. Then, once the database is up and running, the stream platform application will start.

## Consume the application application

You can import and use the Postman collection included among the repository files.

The normal workflow is to use the users API to register as a new user, then authenticate and start using the movies API to create, update, delete and retrieve movies.

Here it follows a description of the implemented endpoints, and how to use them.

### Users API

| URL                  | Method | Usage |
| -------------------- | ------ | ----------- |
| /api/users/register  | POST   | Register a new user |
| /api/users/auth      | POST   | Authenticate an user by generating a JWT token |

### Movies API

| URL                   | Method | Usage     |
| --------------------- | ------ | --------- |
| /api/movies/search    | GET   | Retrieve all the movies present in the patform, according to the given filters
| /api/movies/{id}      | GET   | Retrieve the movie identified by the given identifier |
| /api/movies      | POST   | Create a new movie |
| /api/movies      | PUT   | Update an existing movie |
| /api/movies/{id}      | DELETE   | Delete the movie identified by the given identifier |
