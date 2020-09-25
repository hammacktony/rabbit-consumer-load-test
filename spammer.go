package main

import (
	"flag"
	"io/ioutil"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const exchangeName string = ""

func main() {
	// Get the connection string from the environment variable
	amqpURL := flag.String("url", "amqp://guest:guest@localhost:5672", "AMQP Url")
	workers := flag.Int("workers", 1, "Workers for creating messages")
	messages := flag.Int("messages", 1, "Messages per worker")
	queueName := flag.String("queue", "", "Queue name to send messages too")
	messageLocation := flag.String("filename", "", "Filename of message to spam")

	flag.Parse()

	if *queueName == "" {
		log.Fatal("Need to define a queue to route messages too.")
	}

	if *messageLocation == "" {
		log.Fatal("Need to define a message to spam with.")
	}

	totalMessages := (*workers) * (*messages)

	content, err := ioutil.ReadFile(*messageLocation)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Amqp url is: %s", *amqpURL)
	log.Infof("Total workers: %d", *workers)
	log.Infof("Messages per worker: %d", *messages)
	log.Infof("Total messages being sent overall: %d", totalMessages)

	wg := new(sync.WaitGroup)
	wg.Add(*workers)

	for i := 0; i < *workers; i++ {
		go func(group *sync.WaitGroup) {
			defer group.Done()

			connection, err := amqp.Dial(*amqpURL)
			defer connection.Close()

			if err != nil {
				log.Fatal("Could not establish connection with RabbitMQ:" + err.Error())
			}

			channel, err := connection.Channel()

			if err != nil {
				log.Fatal("Could not open RabbitMQ channel:" + err.Error())
			}

			for j := 0; j < *messages; j++ {
				message := amqp.Publishing{
					Body: content,
				}

				err = channel.Publish(exchangeName, *queueName, false, false, message)

				if err != nil {
					log.Fatal("error publishing a message to the queue:" + err.Error())
				}
			}

		}(wg)
	}
	wg.Wait()
	log.Info("Done spamming")
}
