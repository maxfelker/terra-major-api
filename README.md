Instance API
---

Persistence cloud API to store locations of objects in Unity 3D space

To run the API locally,

```bash
docker compose up --build dev
```

This will run the API locally on [http://localhost:8000](http://localhost:8000). The available endpoints are currently:

```
GET /characters
POST /characters
PATCH /characters/{id}
DELETE /characters/{id}

GET /instances
```

To build the release:

```bash
docker compose up --build release
```

This will run the API locally on [http://localhost](http://localhost). 