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

## Running on Isolated Container
### Requirements
- `docker` installed
- `docker-compose` installed
### Run the image
- Copy `config.example.yml` to `config.yml` and fill in the needed configuration
- run `docker-compose up -d`