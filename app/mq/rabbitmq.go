package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

const DirectExchangeName = "message"

type RabbitMq struct {
	conn *amqp.Connection
}


func (m *RabbitMq) Publish(channel string, p *Payload) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	err = ch.ExchangeDeclare(
		DirectExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	body, _ := p.MarshalBinary()
	err = ch.Publish(
		DirectExchangeName,
		channel,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	return err
}

func (m *RabbitMq) Subscribe(channel string) SubScribeChannel {
	ch, err := m.conn.Channel()
	if err != nil {
		return nil
	}
	err = ch.ExchangeDeclare(
		DirectExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil
	}
	err = ch.QueueBind(
		q.Name,        // queue name
		channel,             // routing key
		DirectExchangeName, // exchange
		false,
		nil)
	msgs, _ := ch.Consume(q.Name,"",true, false, false,false, nil)
	return &RabbitSubscribe{
		channel: msgs,
		ch: ch,
	}
}

func newRabbitMq() MessageQueue {
	link := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		viper.GetString("RabbitMq.User"),
		viper.GetString("RabbitMq.Password"),
		viper.GetString("RabbitMq.Host"),
		viper.GetString("RabbitMq.Port"))
	conn, err := amqp.Dial(link)
	if err != nil {
		panic(fmt.Errorf("rabbitmq error: %w \n", err))
	}
	return &RabbitMq{
		conn: conn,
	}
}

type RabbitSubscribe struct {
	channel  <-chan amqp.Delivery
	ch *amqp.Channel
}

func (subscribe *RabbitSubscribe) Close()  {
	_ = subscribe.ch.Close()
}

func (subscribe *RabbitSubscribe) ReceiveMessage() gjson.Result  {
	msg := <-subscribe.channel
	result := gjson.Parse(string(msg.Body))
	return result
}
