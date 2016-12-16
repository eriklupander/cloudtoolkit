package cloudtoolkit

import (
        "fmt"
        "log"
        "github.com/streadway/amqp"
        "github.com/spf13/viper"
)

type IMessagingClient interface {
        ConnectToBroker() *amqp.Connection
        SendMessage(msg string, contentType string, channel string) error
}

type MessagingClient struct {
        conn *amqp.Connection
}

func (messagingClient *MessagingClient) ConnectToBroker() *amqp.Connection {

        // Replace this ugly stuff with proper props...
        rabbitHost := "192.168.99.100"
        if viper.GetString("profile") != "dev" {
                rabbitHost = "rabbitmq"
        }
        var err error

        messagingClient.conn, err = amqp.Dial("amqp://guest:guest@" + rabbitHost + ":5672/") // 5672
        messagingClient.failOnError(err, "Failed to connect to RabbitMQ")
        // return conn so we can defer close it in the calling func.
        return messagingClient.conn
}

func (messagingClient *MessagingClient) SendMessage(msg string, contentType string, channel string) error {
        span := Tracer.StartSpan("SendMessage")
        defer span.Finish()
        if messagingClient.conn == nil {
                return fmt.Errorf("Connection not initialized")
        }
        ch, err := messagingClient.conn.Channel()
        span.LogEvent("Got channel")
        messagingClient.failOnError(err, "Failed to open a channel")
        defer ch.Close()

        q, err := ch.QueueDeclare(
                channel, // name
                false,   // durable
                false,   // delete when unused
                false,   // exclusive
                false,   // no-wait
                nil,     // arguments
        )
        messagingClient.failOnError(err, "Failed to declare a queue")

        body := msg
        err = ch.Publish(
                "",     // exchange
                q.Name, // routing key
                false,  // mandatory
                false,  // immediate
                amqp.Publishing{
                        ContentType: contentType,
                        Body:        []byte(body),
                })
        Log.Printf(" [x] Sent %s", body)
        messagingClient.failOnError(err, "Failed to publish a message")
        span.LogEvent("Message sent")
        return nil
}


func (c *MessagingClient) failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
                panic(fmt.Sprintf("%s: %s", msg, err))
        }
}
