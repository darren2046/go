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
	Panicerr(err)

	ch, err := conn.Channel()
	Panicerr(err)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	Panicerr(err)

	return &rabbitConnectionStruct{
		connection: conn,
		channel:    ch,
		queue:      q,
	}
}

func (m *rabbitConnectionStruct) Send(data map[string]string) {
	err := m.channel.Publish(
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(map2bin(data)),
		})
	Panicerr(err)
}

func (m *rabbitConnectionStruct) Recv() chan map[string]string {
	msgs, err := m.channel.Consume(
		m.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	Panicerr(err)

	resch := make(chan map[string]string)
	go func() {
		for d := range msgs {
			resch <- bin2map(Str(d.Body))
		}
	}()

	return resch
}
