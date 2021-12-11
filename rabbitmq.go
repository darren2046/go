package golanglibs

import (
	"github.com/streadway/amqp"
)

type rabbitConnectionStruct struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func getRabbitMQ(rabbitMQURL string, queueName string) *rabbitConnectionStruct {
	conn, err := amqp.Dial(rabbitMQURL)
	panicerr(err)

	ch, err := conn.Channel()
	panicerr(err)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return &rabbitConnectionStruct{
		connection: conn,
		channel:    ch,
		queue:      q,
	}
}

func (m *rabbitConnectionStruct) send(data map[string]string) {
	err := m.channel.Publish(
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(map2bin(data)),
		})
	panicerr(err)
}

func (m *rabbitConnectionStruct) recv() chan map[string]string {
	msgs, err := m.channel.Consume(
		m.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	panicerr(err)

	resch := make(chan map[string]string)
	go func() {
		for d := range msgs {
			resch <- bin2map(Str(d.Body))
		}
	}()

	return resch
}
