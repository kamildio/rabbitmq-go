package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Создаём очередь для "info" и "warning"
	queue, err := ch.QueueDeclare(
		"info_warning_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Привязываем очередь к обмену с разными ключами маршрутизации
	severities := []string{"info", "warning"}
	for _, severity := range severities {
		err = ch.QueueBind(
			queue.Name,    // имя очереди
			severity,      // routing key
			"logs_direct", // имя обмена
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to bind a queue: %v", err)
		}
	}

	// Читаем сообщения
	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fmt.Printf("Received: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for info and warning logs...")
	<-forever
}
