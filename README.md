# User service

## Run

make run.local

## Request

#### Local

##### curl

###### User Service

curl -X POST -k http://localhost:8090/v1/user -d '{"user": {"username": "fatboy_slim"}}'
