# MQTT Forwarder to Redis
This app can subscribe to a MQTT topic and forward the received messages to a redis server using Redis Pub/Sub feature

# How to install and Run
## Without Container
### Requirements
- `go` installed
- Redis server installed
### Build and run
- Copy `config.example.yml` to `config.yml` and fill in the needed configuration
- run `go mod download`
- run `go build -o main .`
- run `./main`

## Running on an Isolated Container
### Requirements
- `docker` installed
- `docker-compose` installed
### Run the image
- Copy `config.example.yml` to `config.yml` and fill in the needed configuration
- run `docker-compose up -d`

## Using Redis Outside the Container
### Requirements
- `docker` installed
- `docker-compose` installed
- Redis installed
### Run the image
- Make sure the Redis server is up and running
- Copy `config.example.yml` to `config.yml` and fill in the needed configuration
- [You can access host services using the IP address of the `docker0` interface](https://stackoverflow.com/questions/31324981/how-to-access-host-port-from-docker-container). On linux it is usually `172.17.0.1`. Use `172.17.0.1` as host address on `redis.server` configuration on `config.yml`.
- run `docker-compose up -d app`