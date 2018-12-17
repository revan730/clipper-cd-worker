package queue

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/streadway/amqp"
)

// CDQueue is a name of rabbitMQ queue with CD jobs
const CDQueue = "cdJobs"

// RMQQueue is used for operations with rabbitMQ queues
type RMQQueue struct {
	rabbitConnection *amqp.Connection
	channel          *amqp.Channel
	CDJobsQueue      amqp.Queue
}

func declareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

// NewRMQQueue creates new copy of RMQQueue
func NewRMQQueue(addr string) *RMQQueue {
	conn, err := amqp.Dial(addr)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to rabbitmq: %s", err))
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Couldn't open rabbitmq channel: %s", err))
	}
	cdQueue, err := declareQueue(ch, CDQueue)
	if err != nil {
		panic(fmt.Sprintf("Couldn't declare CD jobs queue: %s", err))
	}
	Queue := &RMQQueue{
		rabbitConnection: conn,
		channel:          ch,
		CDJobsQueue:      cdQueue,
	}

	return Queue
}

// publish sends protobuf encoded message to provided queue
func (q *RMQQueue) publish(msg proto.Message, queue string) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return q.channel.Publish(
		"", queue, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
}

// makeMsgChan creates channel to retranslate message bodies to
func (q *RMQQueue) makeMsgChan(queue string) (chan []byte, error) {
	incomingMsgs, err := q.channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	msgsChan := make(chan []byte)
	go func() {
		for m := range incomingMsgs {
			msgsChan <- m.Body
		}
	}()
	return msgsChan, nil
}

// MakeCDMsgChan creates channel with raw CD job messages
func (q *RMQQueue) MakeCDMsgChan() (chan []byte, error) {
	return q.makeMsgChan(q.CDJobsQueue.Name)
}

// Close gracefully breaks connection with rabbitMQ
func (q *RMQQueue) Close() {
	q.rabbitConnection.Close()
}
