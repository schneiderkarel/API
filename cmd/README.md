# Basic API

Service that enables creating new users into a database, and retrieve them whenever you want to.

## Documentation

* [HTTP API](docs/api.yml)

## Configuration

* `HTTP_SERVER_PORT`: "8080"
* `POSTGRES_DSN`: Database connection string, `required`,
  example: `host=postgres port=5432 user=postgres dbname=postgres sslmode=disable`
  
## Local Development

Local development is possible with or without Docker. In both cases, requirements must be resolved first.

### Run App With Docker

To run application execute the following command (application is automatically rebuild when any source file is
modified):

```
make up
```
* This command creates application environment, and starts the app. 

The service runs on _localhost_.

In order to shut down the service run:

```
make down
```

### Without Docker

To run the app server, or the tests during development use the built-in functions in your IDE with Go installed locally
and environment variables set according to `docker-compose.yml`.
