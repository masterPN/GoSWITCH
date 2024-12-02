## Table of Contents
- [Prerequisites](#prerequisites)
- [Usage](#usage)
  - [Make](#make)
  - [Docker Swarm and Secrets](#docker-swarm-and-secrets)
  - [Environment Variables](#environment-variables)

## Prerequisites
- Install Make (if not already installed)
- Install Docker (if not already installed)

## Usage

### Make

To use Make, run the following command in your terminal:

```bash
make <target>
```

Replace `<target>` with the desired build target, such as `build`, `test`, or `run`.

### Docker Swarm and Secrets

1. Initialize a Docker Swarm:

```bash
docker swarm init
```

This will create a new swarm and make the current node a manager.

2. Create a Docker secret for Redis password:

```bash
echo "your_redis_password" | docker secret create redis_password -
```

Replace `your_redis_password` with your desired password.

### Environment Variables

Create a `.env` file in the `../esl-service/` directory with the following content:

```bash
# ESL-service configuration
SIP_PORT=5060
EXTERNAL_DOMAIN=localhost
```

Create a `.env` file in the `../mssql-service/` directory with the following content:

```bash
# MSSQL-service configuration
PORT=8080
APP_ENV=local

ONEVOIS_DB_HOST=localhost
ONEVOIS_DB_PORT=1433
ONEVOIS_DB_DATABASE=onevois
ONEVOIS_DB_USERNAME=sa
ONEVOIS_DB_PASSWORD=password

WHOLESALE_DB_HOST=localhost
WHOLESALE_DB_PORT=1433
WHOLESALE_DB_DATABASE=wholesale
WHOLESALE_DB_USERNAME=sa
WHOLESALE_DB_PASSWORD=password
```

Create a `.env` file in the `../redis-service/` directory with the following content:

```bash
# Redis-service configuration
PORT=8080
APP_ENV=local

REDIS_HOST=redis:6379
REDIS_PASSWORD=password
```
