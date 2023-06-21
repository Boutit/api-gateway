# Api Gateway

## Run

make run.local

## Request

#### Local

##### curl

###### User Service

curl -X POST -k http://localhost:8090/v1/signup -d '{"email_or_phone": "6472289484"}'

###### Auth Service

curl -X POST -k http://localhost:8090/v1/auth/create_token -d '{"userId": "ey295-asdgjsg-asdgljkas-33dll", "roles": []}'
