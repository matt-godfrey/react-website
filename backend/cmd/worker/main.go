package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/matt-godfrey/react-website/internal/mailer"
	"github.com/matt-godfrey/react-website/internal/queue"
)

func main() {
	err := godotenv.Load()
	ctx := context.Background()

	from := os.Getenv("RESEND_FROM")
	mailerClient := mailer.NewResendMailer(os.Getenv("RESEND_API_KEY"), from)
	if mailerClient == nil {
		log.Fatal("mailerClient is nil")
	}

	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQVHost := os.Getenv("RABBITMQ_VHOST")

	rabbitConn, err := queue.ConnectRabbitMQ(rabbitMQUser, rabbitMQPassword, rabbitMQHost, rabbitMQVHost)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to rabbitmq")

	defer rabbitConn.Close()

	client, err := queue.NewRabbitClient(rabbitConn)
	if err != nil {
		log.Fatal(err)
	}

	messageBus, err := client.Consume("email", "email-service", false)
	if err != nil {
		log.Fatal(err)
	}

	// var blocking chan struct{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		for message := range messageBus {
			log.Printf("Headers: %v", message.Headers)

			email, ok := message.Headers["email"].(string)
			if !ok || email == "" {
				log.Panic("email is empty")
			}

			// Send mail
			msg := mailer.Message{
				// From:    os.Getenv("RESEND_FROM"),
				From:    from,
				To:      []string{message.Headers["email"].(string)},
				Subject: "Login Successful",
				HTML:    "<p>You have successfully logged in to your account.</p>",
				Text:    "You have successfully logged in to your account.",
			}
			_, err = mailerClient.SendMail(ctx, msg)

			if err != nil {
				message.Nack(false, true)
				continue
			}

			if err := message.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
				continue
			}
			log.Printf("Message acknowledged: %s", string(message.Body))
		}
	}()

	log.Println("Consuming; to shutdown press Ctrl+C")
	// <-blocking
	<-ctx.Done()
}
