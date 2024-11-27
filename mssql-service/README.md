# Project mssql-service

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

## VSCode
launch.json
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "cmd/api/main.go",
            "env": {
                "PORT" : "8080",
                "ONEVOIS_DB_HOST" : "",
                "ONEVOIS_DB_PORT" : "1433",
                "ONEVOIS_DB_DATABASE" : "",
                "ONEVOIS_DB_USERNAME" : "",
                "ONEVOIS_DB_PASSWORD" : "",
                "WHOLESALE_DB_HOST" : "",
                "WHOLESALE_DB_PORT" : "1433",
                "WHOLESALE_DB_DATABASE" : "",
                "WHOLESALE_DB_USERNAME" : "",
                "WHOLESALE_DB_PASSWORD" : "",
            }
        }
    ]
}
```

## ENV File
.env at root folder
```properties
PORT=8080
APP_ENV=local

ONEVOIS_DB_HOST=
ONEVOIS_DB_PORT=1433
ONEVOIS_DB_DATABASE=
ONEVOIS_DB_USERNAME=
ONEVOIS_DB_PASSWORD=

WHOLESALE_DB_HOST=
WHOLESALE_DB_PORT=1433
WHOLESALE_DB_DATABASE=
WHOLESALE_DB_USERNAME=
WHOLESALE_DB_PASSWORD=
```
