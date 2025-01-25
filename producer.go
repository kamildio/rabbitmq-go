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

	// Создаём обмен с типом "direct"
	err = ch.ExchangeDeclare(
		"logs_direct", // имя обмена
		"direct",      // тип обмена
		true,          // сохранять ли на диске
		false,         // автудаление
		false,         // эксклюзивный
		false,         // нет ожидания
		nil,           // аргументы
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	// Типы сообщений
	severities := []string{"info", "warning", "error"}

	for _, severity := range severities {
		body := fmt.Sprintf("Log: %s message", severity)
		err = ch.Publish(
			"logs_direct", // обмен
			severity,      // маршрут (routing key)
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		if err != nil {
			log.Fatalf("Failed to publish a message: %v", err)
		}
		fmt.Printf("Sent: %s\n", body)
	}
}
