Instance API
---

Persistence cloud API to store locations of objects in Unity 3D space

To run the API:

```bash
go run cmd/instance-api/main.go
```

This will run the API locally on [http://localhost:8000](http://localhost:8000). The available endpoints are currently:

```
GET /instances
```

To build the binary:

```bash
go build -o bin/instance-api cmd/instance-api/main.go
```

Then run it

```bash
./bin/instance-api
```