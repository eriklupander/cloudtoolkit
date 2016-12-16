package cloudtoolkit

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

const AMQP_CONNECTION_STR string = "amqp.connection.string"

//type IMessagingClient interface {
//        ConnectToBroker() *amqp.Connection
//        SendMessage(msg string, contentType string, channel string) error
//}

type MessagingClient struct {
	conn *amqp.Connection
}

// Typical host amqp://guest:guest@" + rabbitHost + ":5672
func InitMessagingClientFromConfigProperty() *MessagingClient {
	if viper.IsSet(AMQP_CONNECTION_STR) {
		connectionString := viper.GetString(AMQP_CONNECTION_STR)
		return InitMessagingClientFromConnectionString(connectionString)
	} else {
		panic("Config property " + AMQP_CONNECTION_STR + " not set, panicing...")
	}
}

func InitMessagingClientFromConnectionString(connString string) *MessagingClient {
	conn := connectToBroker(connString)
	return &MessagingClient{conn: conn}
}

func connectToBroker(connectionString string) *amqp.Connection {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set. Have you initialized?")
	}

	var err error

	conn, err := amqp.Dial(fmt.Sprintf("%s/", connectionString)) // 5672
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + connectionString)
	}
	// return conn so we can defer close it in the calling func.
	return conn
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
