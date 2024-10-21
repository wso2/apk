# WSO2 APK Code Challenge ToDo Service

## Docker Image

docker pull sampathrajapakse/todo-service:latest

## Docker Image Creation

Use the above provided docker image or you can build using following commands.

docker buildx create --use

docker buildx inspect --bootstrap

docker buildx build --platform linux/amd64,linux/arm64 -t sampathrajapakse/todo-service:latest --push .

## Deploy in K8s Cluster

kubectl create ns apk

kubectl apply -f k8s-artifacts/deployment.yaml -n apk

## Access the pod locally

kubectl port-forward pod/todo-app-<random-id> 8080:8080 -n apk

## Invoke the services

### Retrieve all to-dos

curl http://localhost:8080/todos

### Create a new to-do

curl -X POST http://localhost:8080/todos \
-H "Content-Type: application/json" \
-d '{"task": "Buy groceries", "done": false}'

### Retrieve a single to-do

curl http://localhost:8080/todos/1

### Update an existing to-do

curl -X PUT http://localhost:8080/todos/1 \
-H "Content-Type: application/json" \
-d '{"task": "Buy groceries", "done": true}'

### Delete a to-do

curl -X DELETE http://localhost:8080/todos/1

### Register a user with user-reg header

curl -X POST http://localhost:8080/register \
-H "user-reg: sampath"
