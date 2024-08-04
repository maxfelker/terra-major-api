Instance API
---

Persistence cloud API to store locations of objects in Unity 3D space.

## Run database locally
First, create a `.env` file in the root of the directory with the following values:

```bash
ALLOWED_ORIGINS=http://localhost:5173,https://localhost:5173,http://localhost,https://localhost

POSTGRES_HOST=postgres
POSTGRES_USER=tmdbuser
POSTGRES_PASSWORD=bigland
POSTGRES_DATABASE=terramajor
POSTGRES_SSL=disable

COSMOS_DB_HOST=cosmos
COSMOS_DB_PORT=8081
COSMOS_DB_PRIMARY_KEY=xxx
COSMOS_DB_NAME=terramajor
```

Next, run the local database:

```bash
docker compose up postgres
```

## Run API in Dev mode locally
To run the API locally,

```bash
docker compose up --build dev
```

This will run the API locally on [http://localhost:8000](http://localhost:8000). The available endpoints are currently:

## Run API in Release mode locally 

To build the release:

```bash
docker compose up --build release
```

This will run the API locally on [http://localhost](http://localhost). 