# Apprestaurant_Backend

## Setup Dgraph Database
```
docker run --rm -it -p 8080:8080 -p 9080:9080 -p 8000:8000 -v ~/dgraph:/dgraph dgraph/standalone:v20.03.0
```
**Optional: Add Schema**
```
cd db
curl -X POST localhost:8080/admin/schema --data-binary '@schema.graphql'
```

### Run 
```
go run.
```
