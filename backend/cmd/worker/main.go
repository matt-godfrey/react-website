package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/matt-godfrey/react-website/internal/mailer"
	"github.com/matt-godfrey/react-website/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("worker", "main")
	slog.SetDefault(logger)

	from := os.Getenv("RESEND_FROM")
	mailerClient := mailer.NewResendMailer(os.Getenv("RESEND_API_KEY"), from)
	if mailerClient == nil {
		logger.Warn("mailerClient is nil")
	}

	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQVHost := os.Getenv("RABBITMQ_VHOST")

	rabbitConn, err := queue.ConnectRabbitMQ(rabbitMQUser, rabbitMQPassword, rabbitMQHost, rabbitMQVHost)
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitConn.Close()

	client, err := queue.NewRabbitClient(rabbitConn)
	client.CreateQueue("email", true, false)
	client.CreateExchange("email", amqp.ExchangeDirect, true, false)
	client.CreateBinding("email", "email", "email")
	if err != nil {
		log.Fatal(err)
	}

	messageBus, err := client.Consume("email", "email-service", false)
	if err != nil {
		logger.Warn("Failed to consume email queue", "error", err.Error())
	}

	// var blocking chan struct{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		for message := range messageBus {
			logger.Debug("Headers:", "headers", message.Headers)

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
				HTML:    "<p>You have successfully registered a new account.</p>",
				Text:    "You have successfully registered a new account.",
			}
			_, err = mailerClient.SendMail(ctx, msg)

			if err != nil {
				message.Nack(false, true)
				continue
			}

			if err := message.Ack(false); err != nil {
				logger.Warn("Failed to ack message:", "error", err.Error())
				continue
			}
			logger.Info("Message acknowledged:", "body", string(message.Body))
		}
	}()

	logger.Info("Consuming; to shutdown press Ctrl+C")
	// <-blocking
	<-ctx.Done()
}
