# Serverless GO Crud playground

Go functions interacting with DynamoDB to satisfy CRUD Rest API requests

## Examples

### Get All

```
curl 'https://@functionid.execute-api.eu-central-1.amazonaws.com/dev/animals' \
  -H 'content-type: application/json;charset=UTF-8'
```

### PUT

```
curl -X PUT 'https://@functionid.execute-api.eu-central-1.amazonaws.com/dev/animal' \
  -H 'authorization: Bearer @auth-token' \
  -H 'content-type: application/json;charset=UTF-8' \
  --data-raw '{"id":"004","name":"Fickó","description":"Nagyon csibész kutyus","images":["image1.png","image2.png"]}'
```

### Get

```
curl 'https://@functionid.execute-api.eu-central-1.amazonaws.com/dev/animal?id=003' \
  -H 'content-type: application/json;charset=UTF-8'
```

### Delete

```
curl -X DELETE 'https://@functionid.execute-api.eu-central-1.amazonaws.com/dev/animal?id=003' \
  -H 'authorization: Bearer @auth-token' \
  -H 'content-type: application/json;charset=UTF-8'
```

### Patch

```
curl -X PATCH 'https://@functionid.execute-api.eu-central-1.amazonaws.com/dev/animal?id=003' \
  -H 'authorization: Bearer @auth-token' \
  -H 'content-type: application/json;charset=UTF-8'
  --data-raw '{"status":false}'
```
