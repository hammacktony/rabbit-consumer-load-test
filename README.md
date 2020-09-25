# Spamming tool

Allows for a high velocity and high sustainability load test of messages to a RabbitMQ queue to test one's consumers

## Install

Run the following for dependencies

```sh
go get -u "github.com/sirupsen/logrus"
go get -u "github.com/streadway/amqp"
```

Run `go build spammer.go` to build the binary.

## Run the program

Example:

```sh
./spammer \
    -url=amqp://guest:guest@localhost:5672 \
    -workers=5 \
    -messages=5 \
    -queue=my_events \
    -filename=message.json
```

Sends 5 goroutines to send 5 messages each, so it is 25 messages total.